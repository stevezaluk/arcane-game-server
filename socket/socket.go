package socket

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
)

type Socket struct {
	Listener        *net.Listener
	ConnectionCount int
	MaxConnections  int
	IsClosed        bool
}

func (socket *Socket) Start() {
	uri := "127.0.0.1:" + viper.GetString("port")

	listen, err := net.Listen("tcp", uri)
	if err != nil {
		// log here
		panic(err) // panicing here as this is a fatal error
	}

	socket.Listener = &listen
}

func (socket *Socket) AcceptConnections() {
	for {
		if socket.ConnectionCount == socket.MaxConnections {
			socket.IsClosed = true
			break
		}

		if !socket.IsClosed {
			sock := *socket.Listener
			conn, err := sock.Accept()
			if err != nil {
				// log here
			}

			socket.ConnectionCount += 1
			go socket.HandleClient(conn)
		}
	}
}

func (socket *Socket) HandleClient(conn net.Conn) {
	for {
		buffer := make([]byte, 4096)
		_, err := conn.Read(buffer)
		if err != nil {
			// log here
			// better error checking here
			continue
		}

		fmt.Println("client msg: ", string(buffer))

	}
}

func (socket *Socket) Stop() {
	sock := *socket.Listener
	sock.Close()
}
