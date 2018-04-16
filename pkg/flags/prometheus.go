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

package flags

var (
	// FlagPrometheusServerSecurePort sets the port that Kanali will listen on for incoming requests
	FlagPrometheusServerSecurePort = Flag{
		Long:  "prometheus.secure_port",
		Short: "",
		Value: 0,
		Usage: "Sets the port that Kanali will listen on for incoming requests.",
	}
	// FlagPrometheusServerInsecurePort sets the port that Kanali will listen on for incoming requests
	FlagPrometheusServerInsecurePort = Flag{
		Long:  "prometheus.insecure_port",
		Short: "",
		Value: 9000,
		Usage: "Sets the port that Kanali will listen on for incoming requests.",
	}
	// FlagPrometheusServerInsecureBindAddress specifies the network address that Kanali will listen on for incoming requests
	FlagPrometheusServerInsecureBindAddress = Flag{
		Long:  "prometheus.insecure_bind_address",
		Short: "",
		Value: "0.0.0.0",
		Usage: "Network address that Kanali will listen on for incoming requests.",
	}
	// FlagPrometheusServerSecureBindAddress specifies the network address that Kanali will listen on for incoming requests
	FlagPrometheusServerSecureBindAddress = Flag{
		Long:  "prometheus.secure_bind_address",
		Short: "",
		Value: "0.0.0.0",
		Usage: "Network address that Kanali will listen on for incoming requests.",
	}
)
