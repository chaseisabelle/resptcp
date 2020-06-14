package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/chaseisabelle/resptcp"
	"github.com/tidwall/resp"
)

// example of tcp resp server
func main() {
	// configs
	host := flag.String("host", ":3333", "server host:port")

	flag.Parse()

	// create it
	srv := resptcp.New(*host, handler)

	// listen for server errors
	go func() {
		for err := range srv.Errors {
			println(err.Error())
		}
	}()

	// start the server
	err := srv.Start()

	if err != nil {
		panic(err.Error())
	}
}

// handle incoming connection input and give it a response
func handler(value resp.Value, err error) (resp.Value, error) {
	// handle server read error, if needed
	if err != nil {
		return resp.ErrorValue(err), nil
	}

	// check and see what we got
	println(fmt.Sprintf("%+v", value))

	// build a reply
	switch value.String() {
	case "hello":
		value = resp.SimpleStringValue("hi!")
	case "goodbye":
		value = resp.SimpleStringValue("kthxbye")
	case "integer":
		value = resp.IntegerValue(1)
	case "null":
		value = resp.NullValue()
	default:
		value = resp.ErrorValue(errors.New("wtf?"))
	}

	// respond
	return value, nil
}
