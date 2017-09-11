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
