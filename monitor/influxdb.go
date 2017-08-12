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
	"context"
	"errors"
	"fmt"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/northwesternmutual/kanali/config"
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

	// create db
	q := influx.NewQuery(fmt.Sprintf("CREATE DATABASE %s", viper.GetString(config.FlagInfluxdbDatabase.GetLong())), "", "")
	if response, err := influxClient.Query(q); err != nil || response.Error() != nil {
		return nil, err
	}

	return &InfluxController{Client: influxClient}, nil

}

// WriteRequestData writes contextual request metrics to Influxdb
func (c *InfluxController) WriteRequestData(ctx context.Context) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("influxdb paniced while attempting to writing - this probably means that Kanali was unable to establish a connection on startup")
		}
	}()

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  viper.GetString(config.FlagInfluxdbDatabase.GetLong()),
		Precision: "s",
	})
	if err != nil {
		return err
	}

	pt, err := influx.NewPoint("request_details", getTags(ctx), getFields(ctx), time.Now())
	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	return c.Client.Write(bp)

}

func getTags(ctx context.Context) map[string]string {

	keyName := GetCtxMetric(ctx, "api_key_name")
	if keyName == "" {
		keyName = "none"
	}

	return map[string]string{
		"proxyName":      GetCtxMetric(ctx, "proxy_name"),
		"responseCode":   GetCtxMetric(ctx, "http_response_code"),
		"method":         GetCtxMetric(ctx, "http_method"),
		"keyName":        keyName,
		"proxyNamespace": GetCtxMetric(ctx, "proxy_namespace"),
	}

}

func getFields(ctx context.Context) map[string]interface{} {

	return map[string]interface{}{
		"totalTime":    GetCtxMetric(ctx, "totalTime"),
		"clientIP":     GetCtxMetric(ctx, "client_ip"),
		"responseCode": GetCtxMetric(ctx, "http_response_code"),
		"uri":          GetCtxMetric(ctx, "http_uri"),
	}

}
