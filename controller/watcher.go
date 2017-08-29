// Copyright (c) 2017 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/spec"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	added    = "ADDED"
	modified = "MODIFIED"
	deleted  = "DELETED"
)

// event is an internal struct which we
// we use to hold unmarshalled json events
// from the kubernetes api server
type event struct {
	Type   string
	Object interface{}
}

// rawEvent is an internal struct which we
// we use to hold raw json messages from
// the kubernetes api server
type rawEvent struct {
	Type   string
	Object json.RawMessage
}

// Watch will use goroutines and channels to
// listen to different endpoints on the kubernetes
// api server and act on events that they emit
func (c *Controller) Watch() {

	// make a channel that we will use to move
	// events from one go routine to another
	eventCh := make(chan *event)

	// start a go return that will monitor for
	// new events that are sent through the channel
	go monitor(eventCh, k8sEventHandler{})

	// start listening for events and put
	// them on the channel
	go c.watchResource(eventCh, "apis/kanali.io/v1/apikeies?watch=true")
	go c.watchResource(eventCh, "apis/kanali.io/v1/apikeybindings?watch=true")
	go c.watchResource(eventCh, "apis/kanali.io/v1/apiproxies?watch=true")
	go c.watchResource(eventCh, "api/v1/secrets?fieldSelector=type%3Dkubernetes.io/tls&watch=true")
	go c.watchResource(eventCh, "api/v1/services?watch=true")
	go c.watchResource(eventCh, "api/v1/endpoints?watch=true")

}

func monitor(ev chan *event, h handlerFuncs) {
	for {
		current := <-ev

		switch current.Type {
		case added:
			h.addFunc(current.Object)
		case modified:
			h.updateFunc(current.Object)
		case deleted:
			h.deleteFunc(current.Object)
		}
	}
}

func (c *Controller) watchResource(eventCh chan *event, url string) {
	for {
		if err := c.doWatchResource(eventCh, url); err != nil {
			logrus.Warnf(err.Error())
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *Controller) doWatchResource(eventCh chan *event, url string) error {

	logrus.Infof("attempt to watch %s", url)

	resp, err := c.RestClient.Client.Get(fmt.Sprintf("%s/%s", c.MasterHost, url))
	if err != nil {
		return fmt.Errorf("trouble connecting to k8s apiserver: %s", err.Error())
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf("error closing response body: %s", err.Error())
		}
	}()

	if resp.StatusCode == http.StatusNotFound && strings.Contains(url, "kanali.io") {
		if err := c.CreateTPRs(); err != nil {
			return err
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("k8s apiserver returned a %d status code", resp.StatusCode)
	}

  logrus.Debugf("successfull watch on %s", url)

	decoder := json.NewDecoder(resp.Body)

	for {
		event, err := pollStream(decoder)
		if err == nil {
			eventCh <- event
		} else {
			return err
		}
		if err == io.EOF { // HTTP stream has closed
			return errors.New("http stream to k8s apiserver closed")
		}
	}

}

func pollStream(decoder *json.Decoder) (*event, error) {

	re := &rawEvent{}

	if err := decoder.Decode(re); err != nil {
		return nil, err
	}

	if re.Type == "ERROR" {
		status := &unversioned.Status{}
		if err := json.Unmarshal(re.Object, status); err != nil {
			return nil, fmt.Errorf("error unmarshaling error status: %s", err.Error())
		}
		return nil, fmt.Errorf("received an error event: %s", status.Message)
	}

	event := &event{
		Type: re.Type,
	}

	meta := unversioned.TypeMeta{}

	if err := json.Unmarshal(re.Object, &meta); err != nil {
		return nil, err
	}

	if err := handleValidEvent(meta.Kind, re.Object, event); err != nil {
		return nil, err
	}

	return event, nil

}

func handleValidEvent(kind string, msg json.RawMessage, e *event) error {
	switch kind {
	case "ApiProxy":
		t := spec.APIProxy{}
		if err := json.Unmarshal(msg, &t); err != nil {
			return err
		}
		(*e).Object = t
	case "ApiKeyBinding":
		t := spec.APIKeyBinding{}
		if err := json.Unmarshal(msg, &t); err != nil {
			return err
		}
		(*e).Object = t
	case "ApiKey":
		t := spec.APIKey{}
		if err := json.Unmarshal(msg, &t); err != nil {
			return err
		}
		(*e).Object = t
	case "Secret":
		t := api.Secret{}
		if err := json.Unmarshal(msg, &t); err != nil {
			return err
		}
		(*e).Object = t
	case "Service":
		t := api.Service{}
		if err := json.Unmarshal(msg, &t); err != nil {
			return err
		}
		(*e).Object = t
	case "Endpoints":
		t := api.Endpoints{}
		if err := json.Unmarshal(msg, &t); err != nil {
			return err
		}
		(*e).Object = t
	default:
		return errors.New("unexpected Kubernetes kind")
	}
	return nil
}
