package main

import (
	"flag"
	"fmt"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/resptcp"
)

// example of tcp resp server
func main() {
	// configs
	host := flag.String("host", ":3333", "server host:port")

	flag.Parse()

	// create it
	srv := resptcp.New(*host, handler, '\000')

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
func handler(input []goresp.Value, err error) ([]goresp.Value, error) {
	// define the output
	output := make([]goresp.Value, 0)

	// handle server read error, if needed
	if err != nil {
		return append(output, goresp.NewError(err)), nil
	}

	// inspect what we got
	println(fmt.Sprintf("%+v", input))

	for _, value := range input {
		s, sErr := value.String()
		i, iErr := value.Integer()
		e, eErr := value.Error()
		a, aErr := value.Array()
		f, fErr := value.Float()
		nErr := value.Null()

		println(fmt.Sprintf("%s %+v", s, sErr))
		println(fmt.Sprintf("%d %+v", i, iErr))
		println(fmt.Sprintf("%+v %+v", e, eErr))
		println(fmt.Sprintf("%+v %+v", a, aErr))
		println(fmt.Sprintf("%f %+v", f, fErr))
		println(fmt.Sprintf("%+v", nErr))
	}

	// respond
	return append(output, goresp.NewSimpleString("ok")), nil
}
