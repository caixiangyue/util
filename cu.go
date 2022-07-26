package main

import (
	"flag"
	"fmt"

	"github.com/caixiangyue/util/system"
	"github.com/caixiangyue/util/weibo"
)

var fIP bool
var fWB bool

func init() {
	flag.BoolVar(&fIP, "ip", false, "get ip")
	flag.BoolVar(&fWB, "wb", false, "get ip")
}

func main() {
	flag.Parse()

	if fIP {
		fmt.Print(system.GetPrintedIP())
	} else if fWB {
		weibo.GetTrending()
	} else {
		flag.PrintDefaults()
	}
}
