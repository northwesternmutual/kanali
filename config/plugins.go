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
		FlagPluginsLocation,
		FlagPluginsAPIKeyDecriptionKeyFile,
	)
}

var (
	// FlagPluginsLocation sets the location of custom plugins shared object (.so) files
	FlagPluginsLocation = Flag{
		long:  "plugins.location",
		short: "",
		value: "/",
		usage: "Location of custom plugins shared object (.so) files.",
	}
	// FlagPluginsAPIKeyDecriptionKeyFile set the location of the decryption RSA key file to be used to decrypt incoming API keys.
	FlagPluginsAPIKeyDecriptionKeyFile = Flag{
		long:  "plugin.apiKey.decryption_key_file",
		short: "",
		value: "",
		usage: "Path to valid PEM-encoded private key that matches the public key used to encrypt API keys.",
	}
)
