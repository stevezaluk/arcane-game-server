package server

import (
	"crypto/rsa"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/samber/slog-multi"
	"github.com/spf13/viper"
	"github.com/stevezaluk/arcane-game-server/crypto"
)

type GameServer struct {
	URI             string
	Listener        *net.Listener
	ConnectionCount int
	MaxConnections  int
	IsClosed        bool

	privateKey *rsa.PrivateKey
	publicKey  rsa.PublicKey

	Logger *slog.Logger
}

func (server *GameServer) initLogger() {
	timestamp := time.Now().Format(time.RFC3339Nano)

	filename := viper.GetString("log.path") + "/arcane-" + timestamp + ".json"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	multiHandler := slogmulti.Fanout(
		slog.NewJSONHandler(file, nil),
		slog.NewTextHandler(os.Stdout, nil),
	)

	server.Logger = slog.New(multiHandler)
	slog.SetDefault(server.Logger)
}

func (server *GameServer) initCrypto() {
	slog.Info("Generating RSA-4096 key pair...")
	priv := crypto.GenerateKeyPair()
	server.privateKey = &priv
	server.publicKey = server.privateKey.PublicKey
}

func (server *GameServer) Init() {
	server.initLogger()
	server.initCrypto()

	server.URI = "127.0.0.1:" + viper.GetString("port")
	server.MaxConnections = 8
}

func (server *GameServer) Listen() {
	slog.Info("Starting server...")
	listen, err := net.Listen("tcp", server.URI)
	if err != nil {
		slog.Error("Failed to start listening for connections", "err", err.Error())
		panic(err) // panicing here as this is a fatal error
	}

	slog.Info("Server listening for connections at", "uri", server.URI)
	server.Listener = &listen
}

func (server *GameServer) AcceptConnections() {
	for {
		if server.ConnectionCount == server.MaxConnections {
			server.IsClosed = true
			slog.Warn("Server reached max connection count. No new connections are accepted")
			break
		}

		if !server.IsClosed {
			sock := *server.Listener
			conn, err := sock.Accept()
			if err != nil {
				slog.Error("Failed to accept connection from client ", "client", conn.RemoteAddr().String(), "err", err.Error())
				panic(err)
			}

			slog.Info("Accepted connection from client", "client", conn.RemoteAddr().String())
			server.ConnectionCount += 1
			go server.HandleClient(conn)
		}
	}
}

func (server *GameServer) HandleClient(conn net.Conn) {
	server.NegotiateKeys(conn)
	slog.Info("Client OK. Waiting for JOIN request", "client", conn.RemoteAddr().String())
	for {
		buffer := make([]byte, 6000)
		n, err := conn.Read(buffer)
		if err != nil {
			continue
		}

		plainText := crypto.DecryptMessage(string(buffer[:n]), server.privateKey)
		if plainText == "" {
			conn.Close()
		}
		fmt.Println(plainText)
	}
}

func (server *GameServer) NegotiateKeys(conn net.Conn) {
	slog.Info("Starting key negotiation with client", "client", conn.RemoteAddr().String())

	buffer := make([]byte, 4096)
	n, rErr := conn.Read(buffer)
	if rErr != nil {
		slog.Error("Failed to read buffer from client during key negotiation", "client", conn.RemoteAddr().String())
		conn.Close()
		return
	}

	bufferStr := string(buffer[:n])
	slog.Info("Request from client: ", "buf", bufferStr)
	if bufferStr != "CONNECT" {
		slog.Error("CONNECT Request not formatted properly. Closing connection", "client", conn.RemoteAddr().String())
		conn.Close()
		return
	}

	slog.Info("CONNECT Response acknowledeged. PEM encoding public key...", "client", conn.RemoteAddr().String())

	pemEncodedKey := crypto.PublicKeyToPEM(server.publicKey)
	keyResp := "PUBKEY:" + string(pemEncodedKey)

	slog.Info("Sending public key...")
	_, wErr := conn.Write([]byte(keyResp))
	if wErr != nil {
		slog.Error("Failed to send key to client", "client", conn.RemoteAddr().String())
		conn.Close()
		return
	}
	slog.Info("Key negotiation success for client", "client", conn.RemoteAddr().String())

}

func (server *GameServer) Start() {
	server.Init()
	server.Listen()
	server.AcceptConnections()
}

func (server *GameServer) Stop() {
	sock := *server.Listener
	sock.Close()
}
