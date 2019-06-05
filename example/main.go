package main

import (
	"time"

	variableserver "github.com/dr4ds/go-variable-server"
)

func main() {
	vs := variableserver.New()
	vs.Set("isFalse", true)
	go vs.Start(":9532", "/var")
	n := 0
	for {
		vs.Set("counter", n)
		n++
		time.Sleep(time.Second * 1)
	}
}
