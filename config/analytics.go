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
