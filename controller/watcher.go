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
func (c *Controller) Watch() error {

	// make a channel that we will use to move
	// events from one go routine to another
	eventCh := make(chan *event)

	// start a go return that will monitor for
	// new events that are sent through the channel
	// we will only tolorate one error
	errCh := monitor(eventCh)

	// start listening for events and put
	// them on the channel
	go c.doWatchResource(eventCh, errCh, "apis/kanali.io/v1/apikeies?watch=true")
	go c.doWatchResource(eventCh, errCh, "apis/kanali.io/v1/apikeybindings?watch=true")
	go c.doWatchResource(eventCh, errCh, "apis/kanali.io/v1/apiproxies?watch=true")
	go c.doWatchResource(eventCh, errCh, "api/v1/secrets?fieldSelector=type%3Dkubernetes.io/tls&watch=true")
	go c.doWatchResource(eventCh, errCh, "api/v1/services?watch=true")
	go c.doWatchResource(eventCh, errCh, "api/v1/endpoints?watch=true")

	// if an error is placed on the error channel, return it
	// the reason this works is because listening on either
	// an unbuffered channel or a channel or length one is
	// a BLOCKING operation therefore this function won't
	// return until an error is detected
	return <-errCh

}

// monitor will take action on new events
// that are sent through our event channel
func monitor(ev chan *event) chan error {

	// create buffered channel that we will
	// use to listen for a single error
	errCh := make(chan error, 1)

	// monitor for new events and proccess
	// them base on event type
	go func() {

		// let's start listening forever
		for {

			// a new message has been sent through the channel
			// note that since we are not using a buffered
			// channel, this is a moment in time when both go
			// routines on both sides of the channel are in sync
			current := <-ev

			// let's figure out what type we are dealing with
			// https://github.com/kubernetes/kubernetes/blob/release-1.5.4/pkg/watch/watch.go#L41-L50
			switch current.Type {
			case added:
				doAdd(current.Object)
			case modified:
				logrus.Debugf("udpating... deffering to add")
				doAdd(current.Object)
			case deleted:
				doDelete(current.Object)
			}

		}

	}()

	// return the newly created error channel
	return errCh

}

// doWatchResource is a helper function that will
// perform the watch on our endpoint. When a new event
// is emmited, it will unmarshal it and put it on the channel
func (c *Controller) doWatchResource(eventCh chan *event, errCh chan error, url string) {

	logrus.Infof("attempting to start a watch on %s", url)

	// watch the kubernets api forever
	for {

		// initiate http request
		resp, err := c.RestClient.Client.Get(fmt.Sprintf("%s/%s", c.MasterHost, url))

		if err == nil {
			logrus.Debug("successfully connected to the kubernetes api server without any error")
		}

		// handle errors
		if err != nil {

			logrus.Error(err.Error())

			// we have an error! place it on the channel,
			// clean up, and get out of here
			errCh <- err

			if err := resp.Body.Close(); err != nil {
				logrus.Error(err.Error())
			}

			return

		} else if resp.StatusCode == http.StatusNotFound {

			// even though we are creating the TPRs before this code is executed
			// the Kubernetes server may not finish creating them before this
			// code is executed. Hence, if we get an error indicating that
			// the request resource cannot be found, just let it try again.
			continue

		} else if resp.StatusCode != http.StatusOK {

			// we have an error! place it on the channel,
			// clean up, and get out of here
			errCh <- fmt.Errorf("received a %d status code", resp.StatusCode)

			if err := resp.Body.Close(); err != nil {
				logrus.Error(err.Error())
			}

			return

		}

		// decode request body
		decoder := json.NewDecoder(resp.Body)

		for {

			// we now have an open http stream.
			// listen on it indefinitely
			event, err := pollStream(decoder)

			if err != nil {

				// this is a special kind of error
				// the http stream has closed. This isn't fatal!
				// we just need to reinitiate the connection
				if err == io.EOF {
					break
				}

				// we have an error! place it on the channel,
				// clean up, and get out of here
				errCh <- err

				if err := resp.Body.Close(); err != nil {
					logrus.Error(err.Error())
				}

				return

			}

			// we have a valid new event, pass it along
			eventCh <- event

		}

		// close the response body
		if err := resp.Body.Close(); err != nil {
			logrus.Error(err.Error())
		}

	}

}

// pollStream takes a new event from our http stream
// and process it
func pollStream(decoder *json.Decoder) (*event, error) {

	re := &rawEvent{}
	if err := decoder.Decode(re); err != nil {
		return nil, err
	}

	// alright so you might think the following code is misplaced
	// and should be up above with all of the other event types.
	// However, this we want to be able to put all errors through
	// the error channel and so we'll handle it here
	if re.Type == "ERROR" {

		// error will be of this type
		status := &unversioned.Status{}
		if err := json.Unmarshal(re.Object, status); err != nil {
			// couldn't unmarshall error
			return nil, err
		}
		// return kubernetes error
		return nil, errors.New("kubernetes error")

	}

	// create new event
	event := &event{
		Type: re.Type,
	}

	// we don't know what Kubernetes kind we will unmarshal
	// but we do know that we only want Kubernetes kinds.
	// So let's find out what kind we're dealing with
	meta := unversioned.TypeMeta{}
	if err := json.Unmarshal(re.Object, &meta); err != nil {
		return nil, err
	}

	switch meta.Kind {
	case "ApiProxy":
		t := spec.APIProxy{}
		if err := json.Unmarshal(re.Object, &t); err != nil {
			return nil, err
		}
		event.Object = t
	case "ApiKeyBinding":
		t := spec.APIKeyBinding{}
		if err := json.Unmarshal(re.Object, &t); err != nil {
			return nil, err
		}
		event.Object = t
	case "ApiKey":
		t := spec.APIKey{}
		if err := json.Unmarshal(re.Object, &t); err != nil {
			return nil, err
		}
		event.Object = t
	case "Secret":
		t := api.Secret{}
		if err := json.Unmarshal(re.Object, &t); err != nil {
			return nil, err
		}
		event.Object = t
	case "Service":
		t := api.Service{}
		if err := json.Unmarshal(re.Object, &t); err != nil {
			return nil, err
		}
		event.Object = t
	case "Endpoints":
		t := api.Endpoints{}
		if err := json.Unmarshal(re.Object, &t); err != nil {
			return nil, err
		}
		event.Object = t
	default:
		return nil, errors.New("unexpected Kubernetes kind")
	}

	// return newly created event
	return event, nil

}

// doAdd adds our unmarshalled event to the
// corresponding in memory store
func doAdd(obj interface{}) {

	// use type assertion to determint type and add
	switch obj.(type) {
	case spec.APIProxy:
		if proxy, ok := obj.(spec.APIProxy); ok {
			err := spec.ProxyStore.Set(proxy)
			if err != nil {
				logrus.Errorf("could not add/modify api proxy. skipping: %s", err.Error())
			}
		}
	case spec.APIKey:
		if key, ok := obj.(spec.APIKey); ok {
			if err := key.Decrypt(); err == nil {
				err := spec.KeyStore.Set(key)
				if err != nil {
					logrus.Errorf("could not add/modify apikey. skipping: %s", err.Error())
				}
			} else {
				logrus.Error("could not decrypt apikey. skipping...")
			}
		}
	case spec.APIKeyBinding:
		if binding, ok := obj.(spec.APIKeyBinding); ok {
			err := spec.BindingStore.Set(binding)
			if err != nil {
				logrus.Errorf("could not add/modify apikey binding. skipping: %s", err.Error())
			}
		}
	case api.Secret:
		if secret, ok := obj.(api.Secret); ok {
			err := spec.SecretStore.Set(secret)
			if err != nil {
				logrus.Errorf("could not add/modify secret. skipping: %s", err.Error())
			}
		}
	case api.Service:
		if service, ok := obj.(api.Service); ok {
			err := spec.ServiceStore.Set(spec.CreateService(service))
			if err != nil {
				logrus.Errorf("could not add/modify service. skipping: %s", err.Error())
			}
		}
	case api.Endpoints:
		if endpoints, ok := obj.(api.Endpoints); ok {
			if endpoints.ObjectMeta.Name == "kanali" {
				spec.KanaliEndpoints = &endpoints
			}
		}
	}

}

// doDelete removes our unmarshalled event from the
// corresponding in memory store
func doDelete(obj interface{}) {

	// use type assertion to determint type and delete
	switch obj.(type) {
	case spec.APIProxy:
		if proxy, ok := obj.(spec.APIProxy); ok {
			_, err := spec.ProxyStore.Delete(proxy)
			if err != nil {
				logrus.Errorf("could not delete api proxy. skipping: %s", err.Error())
			}
		}
	case spec.APIKey:
		if key, ok := obj.(spec.APIKey); ok {
			if err := key.Decrypt(); err == nil {
				_, err := spec.KeyStore.Delete(key)
				if err != nil {
					logrus.Errorf("could not delete apikey. skipping: %s", err.Error())
				}
			} else {
				logrus.Warn("could not decrypt apikey. skipping...")
			}
		}
	case spec.APIKeyBinding:
		if binding, ok := obj.(spec.APIKeyBinding); ok {
			_, err := spec.BindingStore.Delete(binding)
			if err != nil {
				logrus.Errorf("could not delete apikey binding. skipping: %s", err.Error())
			}
		}
	case api.Secret:
		if secret, ok := obj.(api.Secret); ok {
			_, err := spec.SecretStore.Delete(secret)
			if err != nil {
				logrus.Errorf("could not delete secret. skipping: %s", err.Error())
			}
		}
	case api.Service:
		if service, ok := obj.(api.Service); ok {
			_, err := spec.ServiceStore.Delete(spec.CreateService(service))
			if err != nil {
				logrus.Errorf("could not delete service. skipping: %s", err.Error())
			}
		}
	}

}
