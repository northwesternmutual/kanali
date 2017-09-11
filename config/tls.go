package config

func init() {
	Flags.Add(
		FlagTLSCertFile,
		FlagTLSKeyFile,
		FlagTLSCaFile,
	)
}

var (
	// FlagTLSCertFile specifies the path to x509 certificate for HTTPS servers.
	FlagTLSCertFile = Flag{
		long:  "tls.cert_file",
		short: "c",
		value: "",
		usage: "Path to x509 certificate for HTTPS servers.",
	}
	// FlagTLSKeyFile pecifies the path to x509 private key matching --tls-cert-file
	FlagTLSKeyFile = Flag{
		long:  "tls.key_file",
		short: "k",
		value: "",
		usage: "Path to x509 private key matching --tls-cert-file.",
	}
	// FlagTLSCaFile specifies the path to x509 certificate authority bundle for mutual TLS
	FlagTLSCaFile = Flag{
		long:  "tls.ca_file",
		short: "",
		value: "",
		usage: "Path to x509 certificate authority bundle for mutual TLS.",
	}
)
