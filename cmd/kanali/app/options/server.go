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

package options

import (
	"github.com/northwesternmutual/kanali/pkg/flags"
)

func init() {
	KanaliOptions.Add(
		FlagServerPort,
		FlagServerBindAddress,
		FlagServerPeerUDPPort,
		FlagServerProxyProtocol,
	)
}

var (
	// FlagServerPort sets the port that Kanali will listen on for incoming requests
	FlagServerPort = flags.Flag{
		Long:  "server.port",
		Short: "p",
		Value: 0,
		Usage: "Sets the port that Kanali will listen on for incoming requests.",
	}
	// FlagServerBindAddress specifies the network address that Kanali will listen on for incoming requests
	FlagServerBindAddress = flags.Flag{
		Long:  "server.bind_address",
		Short: "b",
		Value: "0.0.0.0",
		Usage: "Network address that Kanali will listen on for incoming requests.",
	}
	// FlagServerPeerUDPPort sets the port that all Kanali instances will communicate to each other over
	FlagServerPeerUDPPort = flags.Flag{
		Long:  "server.peer_udp_port",
		Short: "",
		Value: 10001,
		Usage: "Sets the port that all Kanali instances will communicate to each other over.",
	}
	// FlagServerProxyProtocol maintains the integrity of the remote client IP address when incoming traffic to Kanali includes the Proxy Protocol header
	FlagServerProxyProtocol = flags.Flag{
		Long:  "server.proxy_protocol",
		Short: "",
		Value: false,
		Usage: "Maintain the integrity of the remote client IP address when incoming traffic to Kanali includes the Proxy Protocol header.",
	}
)
