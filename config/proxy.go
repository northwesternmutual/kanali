package config

func init() {
	Flags.Add(
		FlagProxyEnableClusterIP,
		FlagProxyHeaderMaskValue,
		FlagProxyEnableMockResponses,
		FlagProxyUpstreamTimeout,
		FlagProxyMaskHeaderKeys,
		FlagProxyTLSCommonNameValidation,
	)
}

var (
	// FlagProxyEnableClusterIP enables to use of cluster ip as opposed to Kubernetes DNS for upstream routing
	FlagProxyEnableClusterIP = Flag{
		long:  "proxy.enable_cluster_ip",
		short: "",
		value: false,
		usage: "Enables to use of cluster ip as opposed to Kubernetes DNS for upstream routing.",
	}
	// FlagProxyHeaderMaskValue sets the value to be used when omitting header values
	FlagProxyHeaderMaskValue = Flag{
		long:  "proxy.header_mask_value",
		short: "",
		value: "ommitted",
		usage: "Sets the value to be used when omitting header values.",
	}
	// FlagProxyEnableMockResponses enables Kanali's mock responses feature. Read the documentation for more information
	FlagProxyEnableMockResponses = Flag{
		long:  "proxy.enable_mock_responses",
		short: "",
		value: false,
		usage: "Enables Kanali's mock responses feature. Read the documentation for more information.",
	}
	// FlagProxyUpstreamTimeout sets the length of upstream timeout
	FlagProxyUpstreamTimeout = Flag{
		long:  "proxy.upstream_timeout",
		short: "",
		value: "0h0m10s",
		usage: "Set length of upstream timeout. Defaults to none",
	}
	// FlagProxyMaskHeaderKeys specifies which headers to mask.
	FlagProxyMaskHeaderKeys = Flag{
		long:  "proxy.mask_header_keys",
		short: "",
		value: []string{},
		usage: "Specify which headers to mask",
	}
	// FlagProxyTLSCommonNameValidation determins whether common name validation occurs as part of an SSL handshake
	FlagProxyTLSCommonNameValidation = Flag{
		long:  "proxy.tls_common_name_validation",
		short: "",
		value: true,
		usage: "Should common name validate as part of an SSL handshake.",
	}
)
