package utils

import (
	"fmt"
	"log"
)

// Panic has bad coping methods for handling errors...
func Panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// LogError logs an error...
func LogError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
