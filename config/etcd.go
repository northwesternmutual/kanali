package config

func init() {
	Flags.Add(
		FlagEtcdCertFile,
		FlagEtcdKeyFile,
		FlagEtcdCaFile,
    FlagEtcdEndpoints,
    FlagEtcdPrefix,
	)
}

var (
	// FlagEtcdCertFile specifies the path to x509 certificate for ETCD servers.
	FlagEtcdCertFile = Flag{
		Long:  "etcd.cert_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 certificate for ETCD servers.",
	}
  // FlagEtcdKeyFile pecifies the path to x509 private key matching --etcd.cert-file
	FlagEtcdKeyFile = Flag{
		Long:  "etcd.key_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 private key matching --etcd.cert_file.",
	}
	// FlagEtcdCaFile specifies the path to x509 certificate authority bundle for mutual TLS
	FlagEtcdCaFile = Flag{
		Long:  "etcd.ca_file",
		Short: "",
		Value: "",
		Usage: "Path to x509 certificate authority bundle for mutual TLS.",
	}
  // FlagEtcdEndpoints specifies a comma delimited list of ETCD hosts.
	FlagEtcdEndpoints = Flag{
		Long:  "etcd.endpoints",
		Short: "",
		Value: []string{},
		Usage: "Comma delimited list of ETCD hosts.",
	}
  // FlagEtcdPrefix specifies the ETCD prefix.
	FlagEtcdPrefix = Flag{
		Long:  "etcd.prefix",
		Short: "",
		Value: "kanali",
		Usage: "ETCD prefix.",
	}
)
