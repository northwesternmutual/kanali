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

package monitor

import (
	"errors"
	"regexp"
	"sync"
	"testing"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	db    string
	store []influx.BatchPoints
	mutex sync.RWMutex
}

func (c *mockClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return 123456789, "", nil
}

func (c *mockClient) Write(bp influx.BatchPoints) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if bp.Database() == "" || c.db != bp.Database() {
		return errors.New("database does not exist")
	}
	c.store = append(c.store, bp)
	return nil
}

func (c *mockClient) Query(q influx.Query) (*influx.Response, error) {
	re := regexp.MustCompile("^CREATE DATABASE (.*)")
	match := re.FindStringSubmatch(q.Command)
	if len(match) != 2 {
		return nil, errors.New("query incorrect")
	}
	if match[1] == "" {
		return nil, errors.New("no database name")
	}
	c.db = match[1]
	return nil, nil
}

func (c *mockClient) Close() error {
	return nil
}

func TestWriteRequestData(t *testing.T) {
	ctlr := &InfluxController{
		Client:    &mockClient{},
		capacity:  viper.GetInt(config.FlagAnalyticsInfluxBufferSize.GetLong()),
		taskQueue: make(chan *influx.Point),
	}
	m := &metrics.Metrics{
		metrics.Metric{Name: "metric-one", Value: "value-one", Index: true},
		metrics.Metric{Name: "metric-two", Value: "value-two", Index: false},
	}

	// not doing anything with items put on this channel
	// but it needs to exist so that we don't block forever
	go func() {
		<-ctlr.taskQueue
	}()

	assert.Nil(t, ctlr.WriteRequestData(m))
}

func TestPrepareWrite(t *testing.T) {
	m := &metrics.Metrics{
		metrics.Metric{Name: "metric-one", Value: "value-one", Index: true},
		metrics.Metric{Name: "metric-two", Value: "value-two", Index: false},
	}

	tags, _ := getTags(m)
	pt, _ := influx.NewPoint(viper.GetString(config.FlagAnalyticsInfluxMeasurement.GetLong()), tags, getFields(m), time.Now())
	bp, _ := prepareWrite([]*influx.Point{pt})

	assert.Equal(t, len(bp.Points()), 1)
}

func TestWrite(t *testing.T) {
	defer viper.Reset()
	ctlr := &InfluxController{Client: &mockClient{}}

	m := &metrics.Metrics{
		metrics.Metric{Name: "metric-one", Value: "value-one", Index: true},
		metrics.Metric{Name: "metric-two", Value: "value-two", Index: false},
	}

	tags, _ := getTags(m)
	pt, _ := influx.NewPoint(viper.GetString(config.FlagAnalyticsInfluxMeasurement.GetLong()), tags, getFields(m), time.Now())
	viper.SetDefault(config.FlagAnalyticsInfluxDb.GetLong(), "test_db")
	bp, _ := prepareWrite([]*influx.Point{pt})
	viper.SetDefault(config.FlagAnalyticsInfluxDb.GetLong(), "")
	assert.Equal(t, ctlr.write(bp).Error(), "no database name")
	viper.SetDefault(config.FlagAnalyticsInfluxDb.GetLong(), "test_db")
	assert.Nil(t, ctlr.write(bp))
	assert.Nil(t, ctlr.write(bp))
	ctlr = &InfluxController{Client: nil}
	assert.Equal(t, ctlr.write(bp).Error(), "influxdb paniced while attempting to write")
}

func TestNewInfluxdbController(t *testing.T) {
	_, err := NewInfluxdbController()
	assert.NotNil(t, err)
	viper.SetDefault(config.FlagAnalyticsInfluxAddr.GetLong(), "http://foo.bar.com")
	_, err = NewInfluxdbController()
	assert.Nil(t, err)
}

func TestRun(t *testing.T) {
	defer viper.Reset()

	client := &mockClient{}
	ctlr := &InfluxController{
		Client:    client,
		capacity:  2,
		taskQueue: make(chan *influx.Point),
	}

	m := &metrics.Metrics{
		metrics.Metric{Name: "metric-one", Value: "value-one", Index: true},
		metrics.Metric{Name: "metric-two", Value: "value-two", Index: false},
	}

	tags, _ := getTags(m)
	pt, _ := influx.NewPoint(viper.GetString(config.FlagAnalyticsInfluxMeasurement.GetLong()), tags, getFields(m), time.Now())
	viper.SetDefault(config.FlagAnalyticsInfluxDb.GetLong(), "test_db")

	go ctlr.Run()

	assert.Equal(t, len(client.store), 0)
	ctlr.taskQueue <- pt
	time.Sleep(1 * time.Millisecond)
	// buffer isn't full, shouldn't have written
	client.mutex.RLock()
	assert.Equal(t, len(client.store), 0)
	client.mutex.RUnlock()
	ctlr.taskQueue <- pt
	time.Sleep(1 * time.Millisecond)
	client.mutex.RLock()
	assert.Equal(t, len(client.store), 1)
	client.mutex.RUnlock()
}

func TestCreateDatabase(t *testing.T) {
	err := createDatabase(&mockClient{})
	assert.Equal(t, err.Error(), "no database name")
	viper.SetDefault(config.FlagAnalyticsInfluxDb.GetLong(), "mydb")
	assert.Nil(t, createDatabase(&mockClient{}))
	viper.SetDefault(config.FlagAnalyticsInfluxDb.GetLong(), "")
}

func TestGetFields(t *testing.T) {
	assert.Equal(t, getFields(&metrics.Metrics{
		metrics.Metric{Name: "metric-one", Value: "value-one", Index: true},
		metrics.Metric{Name: "metric-two", Value: "value-two", Index: false},
	}), map[string]interface{}{
		"metric-one": "value-one",
		"metric-two": "value-two",
	})
}

func TestGetTags(t *testing.T) {
	tags, err := getTags(&metrics.Metrics{
		metrics.Metric{Name: "metric-one", Value: "value-one", Index: true},
		metrics.Metric{Name: "metric-two", Value: "value-two", Index: false},
	})
	assert.Nil(t, err)
	assert.Equal(t, tags, map[string]string{
		"metric-one": "value-one",
	})
	tags, err = getTags(&metrics.Metrics{
		metrics.Metric{Name: "metric-one", Value: 5, Index: true},
		metrics.Metric{Name: "metric-two", Value: "value-two", Index: false},
	})
	assert.Nil(t, tags)
	assert.Equal(t, err.Error(), "InfluxDB requires that the indexed field metric-one be of type string")
}
