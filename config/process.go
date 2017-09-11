package config

func init() {
	Flags.Add(
		FlagProcessLogLevel,
	)
}

var (
	// FlagProcessLogLevel sets the logging level. Choose between 'debug', 'info', 'warn', 'error', 'fatal'
	FlagProcessLogLevel = Flag{
		long:  "process.log_level",
		short: "l",
		value: "info",
		usage: "Sets the logging level. Choose between 'debug', 'info', 'warn', 'error', 'fatal'.",
	}
)
