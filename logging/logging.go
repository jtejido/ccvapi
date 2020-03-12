package logging

import (
	"log"
	"os"
)

type Logging interface {
	Printf(format string, args ...interface{})
}

type serverLogger struct {
	logger *log.Logger
}

func NewLogger(path string) (Logging, error) {
	var lg *log.Logger
	if path == "" {
		lg = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		lg = log.New(f, "", log.LstdFlags)
	}

	return &serverLogger{logger: lg}, nil
}

func (c *serverLogger) Printf(format string, args ...interface{}) {
	c.logger.Printf(format, args...)
}
