package server

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
)

type GameServer struct {
	Listener        *net.Listener
	ConnectionCount int
	MaxConnections  int
	IsClosed        bool
}

func (server *GameServer) Start() {
	uri := "127.0.0.1:" + viper.GetString("port")

	listen, err := net.Listen("tcp", uri)
	if err != nil {
		// log here
		panic(err) // panicing here as this is a fatal error
	}

	server.Listener = &listen
}

func (server *GameServer) AcceptConnections() {
	for {
		if server.ConnectionCount == server.MaxConnections {
			server.IsClosed = true
			break
		}

		if !server.IsClosed {
			sock := *server.Listener
			conn, err := sock.Accept()
			if err != nil {
				// log here
			}

			server.ConnectionCount += 1
			go server.HandleClient(conn)
		}
	}
}

func (server *GameServer) HandleClient(conn net.Conn) {
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

func (server *GameServer) Stop() {
	sock := *server.Listener
	sock.Close()
}
