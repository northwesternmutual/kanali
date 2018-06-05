// Copyright (c) 2018 Northwestern Mutual.
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

var (
	// FlagRSAPublicKeyFile specifies path to RSA public key.
	FlagRSAPublicKeyFile = flags.Flag{
		Long:  "rsa.public_key_file",
		Short: "",
		Value: "",
		Usage: "Path to RSA public key.",
	}
	// FlagKeyName specifies the API key name.
	FlagKeyName = flags.Flag{
		Long:  "key.name",
		Short: "",
		Value: "",
		Usage: "Name of API data.",
	}
	// FlagKeyData specifies existing API key data.
	FlagKeyData = flags.Flag{
		Long:  "key.data",
		Short: "",
		Value: "",
		Usage: "Existing API key data.",
	}
	// FlagKeyOutFile specifies path to RSA public key.
	FlagKeyOutFile = flags.Flag{
		Long:  "key.out_file",
		Short: "",
		Value: "",
		Usage: "Existing API key data.",
	}
	// FlagKeyLength specifies the desired length of the generated API key.
	FlagKeyLength = flags.Flag{
		Long:  "key.length",
		Short: "",
		Value: 32,
		Usage: "Existing API key data.",
	}
)
