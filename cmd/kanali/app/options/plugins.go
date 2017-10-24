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
	flags "github.com/northwesternmutual/kanali/pkg/flags"
)

func init() {
	KanaliOptions.Add(
		FlagPluginsLocation,
		FlagPluginsAPIKeyDecriptionKeyFile,
		FlagPluginsAPIKeyHeaderKey,
	)
}

var (
	// FlagPluginsLocation sets the location of custom plugins shared object (.so) files
	FlagPluginsLocation = flags.Flag{
		Long:  "plugins.location",
		Short: "",
		Value: "/",
		Usage: "Location of custom plugins shared object (.so) files.",
	}
	// FlagPluginsAPIKeyDecriptionKeyFile set the location of the decryption RSA key file to be used to decrypt incoming API keys.
	FlagPluginsAPIKeyDecriptionKeyFile = flags.Flag{
		Long:  "plugins.apiKey.decryption_key_file",
		Short: "",
		Value: "",
		Usage: "Path to valid PEM-encoded private key that matches the public key used to encrypt API keys.",
	}
	// FlagPluginsAPIKeyHeaderKey specifies the name of the HTTP header that will be used to extract the API key.
	FlagPluginsAPIKeyHeaderKey = flags.Flag{
		Long:  "plugins.apiKey.header_key",
		Short: "",
		Value: "apikey",
		Usage: "Name of the HTTP header holding the apikey.",
	}
)
