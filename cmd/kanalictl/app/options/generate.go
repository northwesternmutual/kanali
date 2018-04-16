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
