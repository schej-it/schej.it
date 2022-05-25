package logger

import (
	"io"
	"log"
	"os"
)

var StdOut *log.Logger
var StdErr *log.Logger

func Init(logFile io.Writer) {
	StdOut = log.New(io.MultiWriter(logFile, os.Stdout), "[INFO] ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)
	StdErr = log.New(io.MultiWriter(logFile, os.Stderr), "[ERROR] ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)

	StdOut.Println("######### Server Restarted #########")
}
