package logging

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var Warning = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime)
