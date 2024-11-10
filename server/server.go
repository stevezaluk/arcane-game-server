package server

import (
	"crypto/rsa"
	"fmt"
	"io"
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
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644) // this is not getting closed when the server stops
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

func (server *GameServer) Read(conn net.Conn) (string, error) {
	var ret string

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			slog.Error("Failed to read buffer from client", "err", err.Error(), "client", conn.RemoteAddr().String())
		}
		return ret, err
	}

	ret = string(buffer[:n])

	slog.Debug("Message from client", "msg", ret)
	return ret, nil
}

func (server *GameServer) ReadEncrypted(conn net.Conn) (string, error) {
	return "", nil
}

func (server *GameServer) HandleClient(conn net.Conn) {
	if !server.NegotiateKeys(conn) {
		slog.Error("Key negotiation failed for client. Closing connection...")
		conn.Close()
		return
	}

	slog.Info("Key negotiation success for client", "client", conn.RemoteAddr().String())
	slog.Info("Client OK. Waiting for JOIN request", "client", conn.RemoteAddr().String())
	for {
		cipherText, err := server.Read(conn)
		if err != nil {
			if err == io.EOF {
				slog.Info("Client has disconnected", "client", conn.RemoteAddr().String())
			}
			conn.Close()
			break
		}

		plainText := crypto.DecryptMessage(cipherText, server.privateKey)
		if plainText == "" {
			conn.Close()
			break
		}
		fmt.Println(plainText)
	}
}

func (server *GameServer) NegotiateKeys(conn net.Conn) bool {
	slog.Info("Starting key negotiation with client", "client", conn.RemoteAddr().String())

	var result bool

	buffer, err := server.Read(conn)
	if err != nil {
		return result
	}

	if buffer != "CONNECT" {
		slog.Error("CONNECT Request not formatted properly. Closing connection", "client", conn.RemoteAddr().String())
		return result
	}

	slog.Info("CONNECT Response acknowledeged. PEM encoding public key...", "client", conn.RemoteAddr().String())

	pemEncodedKey := crypto.PublicKeyToPEM(server.publicKey)
	keyResp := "PUBKEY:" + string(pemEncodedKey)

	slog.Info("Sending public key...")
	_, wErr := conn.Write([]byte(keyResp))
	if wErr != nil {
		slog.Error("Failed to send key to client", "client", conn.RemoteAddr().String())
		return result
	}

	result = true
	return result
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
