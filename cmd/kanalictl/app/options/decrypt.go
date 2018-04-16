package options

import (
	"github.com/northwesternmutual/kanali/pkg/flags"
)

var (
	// FlagRSAPrivateKeyFile specifies path to RSA private key.
	FlagRSAPrivateKeyFile = flags.Flag{
		Long:  "rsa.private_key_file",
		Short: "",
		Value: "",
		Usage: "Path to RSA private key.",
	}
	// FlagKeyInFile specifies path to file or directory containing API key resources.
	FlagKeyInFile = flags.Flag{
		Long:  "key.in_file",
		Short: "",
		Value: "",
		Usage: "Path to file or directory containing API key resources",
	}
)
