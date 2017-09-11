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
	)
}

var (
	// FlagAnalyticsInfluxAddr specifies the Influxdb address. Address should be of the form 'http://host:port' or 'http://[ipv6-host%zone]:port'
	FlagAnalyticsInfluxAddr = Flag{
		long:  "analytics.influx_addr",
		short: "",
		value: "http://monitoring-influxdb.kube-system.svc.cluster.local:8086",
		usage: "InfluxDB address. Address should be of the form 'http://host:port' or 'http://[ipv6-host%zone]:port'.",
	}
	// FlagAnalyticsInfluxDb specifies the name of the InfluxDB database
	FlagAnalyticsInfluxDb = Flag{
		long:  "analytics.influx_db",
		short: "",
		value: "k8s",
		usage: "InfluxDB database name",
	}
	// FlagAnalyticsInfluxUsername specifies the InfluxDB username
	FlagAnalyticsInfluxUsername = Flag{
		long:  "analytics.influx_username",
		short: "",
		value: "",
		usage: "InfluxDB username",
	}
	// FlagAnalyticsInfluxPassword specifies the InfluxDB password
	FlagAnalyticsInfluxPassword = Flag{
		long:  "analytics.influx_password",
		short: "",
		value: "",
		usage: "InfluxDB password",
	}
)
