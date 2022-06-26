package logger

import "log"

func Catch(err error) {
	if err != nil {
		log.Panic(err)
	}
}
