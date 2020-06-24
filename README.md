# RESPTCP
a simple, lightweight TCP server package for building simple RESP apps

---

TCP: https://en.wikipedia.org/wiki/Transmission_Control_Protocol

RESP: https://redis.io/topics/protocol

---

### how it works

1. client sends a single RESP message (i.e. `+hello\r\n`)
2. server receives single RESP message over TCP
3. server passes single RESP message to handler
4. handler return RESP message(s) to server
5. server responds to client with RESP message(s)
6. client does whatever

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

logs...
```
hello
goodbye
integer
null
[poop]
```

---

### limitations

the biggest limitation is that the server only handles 
a single incoming message at a time. if you want to
handle more than one line, you'll need to send it as
an array of lines.

for example, if you send `+hello\r\n+goodbye\r\n`, the
server will call the handler twice, the first time to
handle `+hello\r\n`, and the second time to handle
`+goodbye\r\n`.

in the example above, if you wanted to handle both lines
in the same handler call, you would need to send it as
an array: `*2\r\n+hello\r\n+goodbye\r\n`.