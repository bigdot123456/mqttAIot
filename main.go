package main

import (
	_ "client/boot"
	_ "client/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
