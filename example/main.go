package main

import (
	"flag"
	"fmt"
	"github.com/chaseisabelle/resptcp"
	"github.com/tidwall/resp"
)

func main() {
	host := flag.String("host", ":3333", "server host:port")

	flag.Parse()

	err := resptcp.New(*host, handler).Start()

	if err != nil {
		panic(err.Error())
	}
}

func handler(value resp.Value, err error) (resp.Value, error) {
	println(fmt.Sprintf("type: %+v\nvalue: %+v", value.Type(), value))

	return value, nil
}
