package socket

import (
	"net"

	"github.com/spf13/viper"
)

type Socket struct {
	Listener *net.Listener
}

func (socket *Socket) Start() {
	uri := "localhost:" + viper.GetString("port")

	listen, err := net.Listen("tcp", uri)
	if err != nil {
		// log here
		panic(err) // panicing here as this is a fatal error
	}

	socket.Listener = &listen
}

func (socket *Socket) WaitForConnections() {

}
