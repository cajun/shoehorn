package logger

import (
	"fmt"
)

type Status struct {
	Done    bool
	Message string
}

var (
	communication = make(chan Status)
	serverOn      bool
)

func New(serverOn bool) chan Status {
	serverOn = serverOn
	return communication
}

func Log(message string) {
	communication <- Status{Message: message}
}

func Done() {
	communication <- Status{Done: true}
}

func InitStatus() Status {
	return Status{Done: false}
}

func Write(status Status) {
	if serverOn {
		// do server stuff
	} else {
		fmt.Println(status.Message)
	}
}
