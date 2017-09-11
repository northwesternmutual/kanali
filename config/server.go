package config

func init() {
	Flags.Add(
		FlagServerPort,
		FlagServerBindAddress,
		FlagServerPeerUDPPort,
		FlagServerProxyProtocol,
	)
}

var (
	// FlagServerPort sets the port that Kanali will listen on for incoming requests
	FlagServerPort = Flag{
		long:  "server.port",
		short: "p",
		value: 0,
		usage: "Sets the port that Kanali will listen on for incoming requests.",
	}
	// FlagServerBindAddress specifies the network address that Kanali will listen on for incoming requests
	FlagServerBindAddress = Flag{
		long:  "server.bind_address",
		short: "b",
		value: "0.0.0.0",
		usage: "Network address that Kanali will listen on for incoming requests.",
	}
	// FlagServerPeerUDPPort sets the port that all Kanali instances will communicate to each other over
	FlagServerPeerUDPPort = Flag{
		long:  "server.peer_udp_port",
		short: "",
		value: 10001,
		usage: "Sets the port that all Kanali instances will communicate to each other over.",
	}
	// FlagServerProxyProtocol maintains the integrity of the remote client IP address when incoming traffic to Kanali includes the Proxy Protocol header
	FlagServerProxyProtocol = Flag{
		long:  "server.proxy_protocol",
		short: "",
		value: false,
		usage: "Maintain the integrity of the remote client IP address when incoming traffic to Kanali includes the Proxy Protocol header.",
	}
)
