package logging

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
var Status = log.New(os.Stdout, "STATUS: ", log.Ldate|log.Ltime)
var Warning = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime)
