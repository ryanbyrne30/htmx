package main

import (
	"log"
	"os"
)

// Custom leveled info logger
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Custom leveled error logger
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)