# RESPTCP
a simple, lightweight TCP server package for building RESP apps

---

TCP: https://en.wikipedia.org/wiki/Transmission_Control_Protocol

RESP: https://redis.io/topics/protocol

---

### example

`/example/main.go`
```go
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
```

running the example:
```
$ printf "+hello\r\n+goodbye\r\n+integer\r\n+null\r\npoop\r\n" | netcat 127.0.0.1 3333
+hi!
+kthxbye
:1
$-1
-wtf?
```
