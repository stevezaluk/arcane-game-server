package server

import (
	"crypto/rsa"
	"net"
	"strings"

	"github.com/spf13/viper"
	"github.com/stevezaluk/arcane-game-server/crypto"
)

type GameServer struct {
	Listener        *net.Listener
	ConnectionCount int
	MaxConnections  int
	IsClosed        bool

	privateKey rsa.PrivateKey
	publicKey  rsa.PublicKey
}

func (server *GameServer) Start() {
	priv, pub := crypto.GenerateKeyPair()
	server.privateKey = priv
	server.publicKey = pub

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
			go server.NegotiateKeys(conn)
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
	}
}

func (server *GameServer) NegotiateKeys(conn net.Conn) {
	negotiationSuccess := false

	buffer := make([]byte, 4096)
	_, err := conn.Read(buffer)
	if err != nil {
		panic(err) // transmission error during key negotation
	}

	bufferStr := string(buffer)
	if strings.HasPrefix(bufferStr, "JOIN:") {
		keyResp := "PUBKEY:" + string(crypto.PublicKeyToPEM(server.publicKey))
		_, err := conn.Write([]byte(keyResp))
		if err != nil {
			panic(err) // transmission error during key negotiation
		}
		negotiationSuccess = true
	}

	if negotiationSuccess {
		server.HandleClient(conn)
	} else {
		conn.Close()
	}

}

func (server *GameServer) Stop() {
	sock := *server.Listener
	sock.Close()
}
