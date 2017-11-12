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

package config

func init() {
	Flags.Add(
		FlagAnalyticsInfluxAddr,
		FlagAnalyticsInfluxDb,
		FlagAnalyticsInfluxUsername,
		FlagAnalyticsInfluxPassword,
		FlagAnalyticsInfluxBufferSize,
		FlagAnalyticsInfluxMeasurement,
	)
}

var (
	// FlagAnalyticsInfluxAddr specifies the Influxdb address. Address should be of the form 'http://host:port' or 'http://[ipv6-host%zone]:port'
	FlagAnalyticsInfluxAddr = Flag{
		Long:  "analytics.influx_addr",
		Short: "",
		Value: "http://monitoring-influxdb.kube-system.svc.cluster.local:8086",
		Usage: "InfluxDB address. Address should be of the form 'http://host:port' or 'http://[ipv6-host%zone]:port'.",
	}
	// FlagAnalyticsInfluxDb specifies the name of the InfluxDB database
	FlagAnalyticsInfluxDb = Flag{
		Long:  "analytics.influx_db",
		Short: "",
		Value: "k8s",
		Usage: "InfluxDB database name",
	}
	// FlagAnalyticsInfluxUsername specifies the InfluxDB username
	FlagAnalyticsInfluxUsername = Flag{
		Long:  "analytics.influx_username",
		Short: "",
		Value: "",
		Usage: "InfluxDB username",
	}
	// FlagAnalyticsInfluxPassword specifies the InfluxDB password
	FlagAnalyticsInfluxPassword = Flag{
		Long:  "analytics.influx_password",
		Short: "",
		Value: "",
		Usage: "InfluxDB password",
	}
	// FlagAnalyticsInfluxBufferSize specifies the InfluxDB buffer size.
	// Request metrics will be written to InfluxDB when this buffer is full.
	FlagAnalyticsInfluxBufferSize = Flag{
		Long:  "analytics.influx_buffer_size",
		Short: "",
		Value: 10,
		Usage: "InfluxDB buffer size. Request metrics will be written to InfluxDB when this buffer is full.",
	}
	// FlagAnalyticsInfluxMeasurement specifies the InfluxDB measurement to be used for Kanali request metrics.
	FlagAnalyticsInfluxMeasurement = Flag{
		Long:  "analytics.influx_measurement",
		Short: "",
		Value: "request_details",
		Usage: " InfluxDB measurement to be used for Kanali request metrics.",
	}
)
