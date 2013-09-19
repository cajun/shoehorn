package logger

import (
	"fmt"
	"io"
	"time"
)

type Status struct {
	Done    bool
	message string
}

var (
	communication = make(chan Status)
	output        io.Writer
)

func New(w io.Writer) chan Status {
	SetWriter(w)
	return communication
}

func SetWriter(w io.Writer) {
	output = w
}

func Log(message string) {
	communication <- Status{message: message}
}

func Done() {
	communication <- Status{Done: true}
}

func InitStatus() Status {
	return Status{Done: false}
}

func Write(status Status) {
	const layout = "Jan 2, 2006 at 3:04.05 pm (MST)"
	t := time.Now().Format(layout)
	output.Write([]byte(fmt.Sprintf("%s -- %s", t, status.message)))
}
