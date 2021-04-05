package configs

import (
	"log"
	"os"
)

func NewLog() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}
