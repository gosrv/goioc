package util

import (
	"fmt"
	"github.com/xiaotie/gcluster/util"
	"log"
)

func Assert(condition bool, msg string) {
	if !condition {
		Panic(msg)
	}
}

func VerifyNotNull(ins interface{}) {
	if ins == nil {
		util.Panic("nil interface")
	}
}

func Panic(format string, a ...interface{}) {
	log.Panic(fmt.Sprintf(format, a...))
}
