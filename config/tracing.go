package config

func init() {
	Flags.Add(
		FlagTracingJaegerServerURL,
		FlagTracingJaegerAgentURL,
	)
}

var (
	// FlagTracingJaegerServerURL specifies the endpoint to the Jaeger server
	FlagTracingJaegerServerURL = Flag{
		long:  "tracing.jaeger_server_url",
		short: "",
		value: "jaeger-all-in-one-agent.default.svc.cluster.local",
		usage: "Endpoint to the Jaeger server",
	}
	// FlagTracingJaegerAgentURL specifies the endpoint to the Jaeger agent
	FlagTracingJaegerAgentURL = Flag{
		long:  "tracing.jaeger_agent_url",
		short: "",
		value: "jaeger-all-in-one-agent.default.svc.cluster.local",
		usage: "Endpoint to the Jaeger agent",
	}
)
