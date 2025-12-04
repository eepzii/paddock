package app

import (
	"log"
	"time"
)

func PrintSuccess(token string, startTime time.Time) {
	msg := Message{
		Token:     token,
		Success:   true,
		StartTime: startTime,
		Error:     nil,
	}
	msgJSON, err := msg.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(msgJSON))
}

func PrintFatal(err error, startTime time.Time) {
	msg := Message{
		Token:     "",
		Success:   false,
		StartTime: startTime,
		Error:     err,
	}
	msgJSON, err := msg.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(string(msgJSON))
}
