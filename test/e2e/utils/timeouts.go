package utils

import "time"

var (
	Poll                     = 2 * time.Second
	NamespaceCreationTimeout = 30 * time.Second
	NamespaceCleanupTimeout  = 15 * time.Minute
)
