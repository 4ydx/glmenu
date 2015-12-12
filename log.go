package glmenu

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var IsDebug = false

type MenuLogger struct {
	*log.Logger
}

func NewMenuLogger(location string) (*MenuLogger, error) {
	f, err := os.OpenFile(location, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	mn := &MenuLogger{
		Logger: log.New(f, "", log.Lshortfile),
	}
	return mn, nil
}

func (nm *MenuLogger) Close() error {
	return nm.Close()
}

func MenuDebug(message string) {
	if IsDebug {
		pc, fn, line, _ := runtime.Caller(1)
		fmt.Printf("[error] in %s[%s:%d] %s", runtime.FuncForPC(pc).Name(), fn, line, message)
	}
}
