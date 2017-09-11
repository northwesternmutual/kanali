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
		Long:  "proxy.enable_cluster_ip",
		Short: "",
		Value: false,
		Usage: "Enables to use of cluster ip as opposed to Kubernetes DNS for upstream routing.",
	}
	// FlagProxyHeaderMaskValue sets the Value to be used when omitting header Values
	FlagProxyHeaderMaskValue = Flag{
		Long:  "proxy.header_mask_Value",
		Short: "",
		Value: "ommitted",
		Usage: "Sets the Value to be used when omitting header Values.",
	}
	// FlagProxyEnableMockResponses enables Kanali's mock responses feature. Read the documentation for more information
	FlagProxyEnableMockResponses = Flag{
		Long:  "proxy.enable_mock_responses",
		Short: "",
		Value: false,
		Usage: "Enables Kanali's mock responses feature. Read the documentation for more information.",
	}
	// FlagProxyUpstreamTimeout sets the length of upstream timeout
	FlagProxyUpstreamTimeout = Flag{
		Long:  "proxy.upstream_timeout",
		Short: "",
		Value: "0h0m10s",
		Usage: "Set length of upstream timeout. Defaults to none",
	}
	// FlagProxyMaskHeaderKeys specifies which headers to mask.
	FlagProxyMaskHeaderKeys = Flag{
		Long:  "proxy.mask_header_keys",
		Short: "",
		Value: []string{},
		Usage: "Specify which headers to mask",
	}
	// FlagProxyTLSCommonNameValidation determins whether common name validation occurs as part of an SSL handshake
	FlagProxyTLSCommonNameValidation = Flag{
		Long:  "proxy.tls_common_name_validation",
		Short: "",
		Value: true,
		Usage: "Should common name validate as part of an SSL handshake.",
	}
)
