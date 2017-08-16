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
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/spf13/viper"
)

// InfluxController represents configuration to create an Influxdb connection
type InfluxController struct {
	Client influx.Client
}

// NewInfluxdbController creates a new controller allowing
// connection to Influxdb
func NewInfluxdbController() (*InfluxController, error) {
	influxClient, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     viper.GetString(config.FlagInfluxdbAddr.GetLong()),
		Username: viper.GetString(config.FlagInfluxdbUsername.GetLong()),
		Password: viper.GetString(config.FlagInfluxdbPassword.GetLong()),
	})
	if err != nil {
		return nil, err
	}
	return &InfluxController{Client: influxClient}, nil
}

// WriteRequestData writes contextual request metrics to Influxdb
func (c *InfluxController) WriteRequestData(m *metrics.Metrics) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("influxdb paniced while attempting to write")
		}
	}()
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: viper.GetString(config.FlagInfluxdbDatabase.GetLong()),
	})
	if err != nil {
		return err
	}
	tags, err := getTags(m)
	if err != nil {
		return err
	}
	pt, err := influx.NewPoint("request_details", tags, getFields(m), time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)
	if err := c.Client.Write(bp); err == nil {
		return nil
	}
	if err := createDatabase(c.Client); err != nil {
		return err
	}
	return c.Client.Write(bp)
}

func createDatabase(c influx.Client) error {
	q := influx.NewQuery(fmt.Sprintf("CREATE DATABASE %s", viper.GetString(config.FlagInfluxdbDatabase.GetLong())), "", "")
	response, err := c.Query(q)
	if err != nil {
		return err
	}
	if response != nil && response.Error() != nil {
		return response.Error()
	}
	return nil
}

func getTags(m *metrics.Metrics) (map[string]string, error) {
	tags := make(map[string]string)
	for _, metric := range *m {
		if !metric.Index {
			continue
		}
		tagValue, ok := metric.Value.(string)
		if !ok {
			return nil, fmt.Errorf("InfluxDB requires that the indexed field %s be of type string", metric.Name)
		}
		tags[metric.Name] = tagValue
	}
	return tags, nil
}

func getFields(m *metrics.Metrics) map[string]interface{} {
	fields := make(map[string]interface{})
	for _, metric := range *m {
		fields[metric.Name] = metric.Value
	}
	return fields
}
