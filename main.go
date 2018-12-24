package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	argsNum := len(os.Args)
	if argsNum < 2 {
		report("缺少参数：目的网址或ip")
		return
	}
	var (
		host     = os.Args[1]
		loop     = 4
		dataSize = 64
	)
	ping := &pinger{
		Runloop: loop,
		Timeout: time.Duration(400) * time.Millisecond,
		Data:    make([]byte, dataSize),
	}
	ping.ping(host)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func report(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprint(os.Stdout, "\n")
}
