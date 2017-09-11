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
