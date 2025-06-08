package logger

import (
	"log"
	"os"
)

var STREAM = os.Stdout
var INFOLOG = log.New(STREAM, "[INFO] ", log.Ldate|log.Ltime)
var WARNINGLOG = log.New(STREAM, "[WARNING] ", log.Ldate|log.Ltime)
var ERRORLOG = log.New(STREAM, "[ERROR] ", log.Ldate|log.Ltime)
