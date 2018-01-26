package server

import (
	"os"
)

var shutdownSignals = []os.Signal{os.Interrupt}
