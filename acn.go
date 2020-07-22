package acn

import (
	"log"
)

// Config is the Navigator config
type Config struct {
	// Folder is the folder of the Go project
	Folder string
}

// Navigator is a stateful code navigator that will feed matching AST lines to the caller
type Navigator struct {
	logger *log.Logger
	conf   Config
}

func NewNavigator(logger *log.Logger, conf Config) Navigator {
	return Navigator{
		logger: logger,
		conf:   conf,
	}
}
