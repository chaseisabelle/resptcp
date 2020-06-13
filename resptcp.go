package resptcp

import (
	"bufio"
	"github.com/chaseisabelle/stop"
	"github.com/tidwall/resp"
	"io"
	"net"
)

type Server struct {
	host    string
	handler func(resp.Value, error) (resp.Value, error)
	Errors  chan error
}

func New(host string, handler func(resp.Value, error) (resp.Value, error)) *Server {
	return &Server{
		host:    host,
		handler: handler,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp4", s.host)

	if err != nil {
		return err
	}

	defer func() {
		err = listener.Close()

		if err != nil {
			panic(err.Error())
		}
	}()

	for !stop.Stopped() {
		connection, err := listener.Accept()

		if err != nil {
			return err
		}

		go func() {
			reader := resp.NewReader(bufio.NewReader(connection))

			for {
				value, _, err := reader.ReadValue()

				if err == io.EOF {
					return
				}

				value, err = s.handler(value, err)

				if err != nil {
					s.Errors <- err

					continue
				}

				marshalled, err := value.MarshalRESP()

				if err != nil {
					s.Errors <- err

					continue
				}

				_, err = connection.Write(marshalled)

				if err != nil {
					s.Errors <- err

					continue
				}
			}
		}()
	}
}
