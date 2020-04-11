package error

import (
	"fmt"
	"log"
	"runtime/debug"
)

type MyError struct {
	Err        error
	Message    string
	StackTrace string
	Mics       map[string]string
}

func (e MyError) Error() string {
	return e.Message
}

func WrapError(err error, messagef string, args ...interface{}) MyError {
	return MyError{
		Err:        err,
		Message:    fmt.Sprintf(messagef, args...),
		StackTrace: string(debug.Stack()),
		Mics:       make(map[string]string),
	}
}

func HandleError(logId int, err error) {
	log.SetPrefix(fmt.Sprintf("[logId: %d]", logId))
	log.Printf("%#v", err)
	fmt.Printf("[logId: %d] %s", logId, err.Error())
}
