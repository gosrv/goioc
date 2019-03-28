package util

import (
	"fmt"
	"log"
)

func Assert(condition bool, msg string) {
	if !condition {
		Panic(msg)
	}
}

func VerifyNotNull(ins interface{}) {
	if ins == nil {
		Panic("nil interface")
	}
}

func VerifyNoError(err error) {
	if err != nil {
		Panic(err.Error())
	}
}

func Panic(format string, a ...interface{}) {
	log.Panic(fmt.Sprintf(format, a...))
}
