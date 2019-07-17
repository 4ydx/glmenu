package glmenu

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// IsDebug when set to true outputs debug logging information
var IsDebug = false

func NewMenuLogger(location string) (*log.Logger, error) {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	mn := log.New(f, "", log.Lshortfile)
	return mn, nil
}

func MenuDebug(message string) {
	if IsDebug {
		pc, fn, line, _ := runtime.Caller(1)
		fmt.Printf("[error] in %s[%s:%d] %s", runtime.FuncForPC(pc).Name(), fn, line, message)
	}
}
