package logger

import (
	"os"
	"log"
	"io/ioutil"
)

type LogLevel int

const (
	INFO  LogLevel = iota
	DEBUG
)

var logLevelByName = map[string]LogLevel{
	"INFO":  INFO,
	"DEBUG": DEBUG,
}

var logLevel LogLevel = INFO

var Info *log.Logger
var Debug *log.Logger
var Return *log.Logger

// creats the Info and Debug logger
// Debug logger logs only if the env variable LOGLEVEL is set to DEBUG
// (INFO is the default value)
func init() {
	if val, ok := logLevelByName[os.Getenv("LOGLEVEL")]; ok {
		logLevel = val
	}
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	if logLevel == DEBUG {
		Debug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
	} else {
		Debug = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime)
	}
	Return = log.New(os.Stdout, "", 0)
}

// checks if err is != null and calls panic if it is
func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
