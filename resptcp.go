package resptcp

import (
	"bufio"
	"github.com/chaseisabelle/goresp"
	"github.com/chaseisabelle/stop"
	"io"
	"net"
)

type Server struct {
	host      string
	handler   func([]goresp.Value, error) ([]goresp.Value, error)
	delimiter byte
	Errors    chan error
}

func New(host string, handler func([]goresp.Value, error) ([]goresp.Value, error), delimiter byte) *Server {
	return &Server{
		host:      host,
		handler:   handler,
		delimiter: delimiter,
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
			s.Errors <- err
		}
	}()

	for !stop.Stopped() {
		connection, err := listener.Accept()

		if err != nil {
			return err
		}

		go func() {
			reader := bufio.NewReader(connection)

			for {
				input, err := reader.ReadBytes(s.delimiter)

				if err == io.EOF {
					return
				}

				if err != nil {
					s.Errors <- err

					continue
				}

				values, err := goresp.Decode(input[:len(input)-1])
				values, err = s.handler(values, err)

				if err != nil {
					s.Errors <- err

					continue
				}

				output, err := goresp.Encode(values)

				if err != nil {
					s.Errors <- err

					continue
				}

				_, err = connection.Write(output)

				if err != nil {
					s.Errors <- err
				}
			}
		}()
	}

	return nil
}
