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
	arcaneErrors "github.com/stevezaluk/arcane-game-server/errors"
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

func (server *GameServer) initLogger() error {
	timestamp := time.Now().Format(time.RFC3339Nano)

	filename := viper.GetString("log.path") + "/arcane-" + timestamp + ".json"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644) // this is not getting closed when the server stops
	if err != nil {
		return arcaneErrors.ErrLogFileFailed
	}

	multiHandler := slogmulti.Fanout(
		slog.NewJSONHandler(file, nil),
		slog.NewTextHandler(os.Stdout, nil),
	)

	server.Logger = slog.New(multiHandler)
	slog.SetDefault(server.Logger)

	return nil
}

func (server *GameServer) initCrypto() error {
	slog.Info("Generating RSA-4096 key pair...")
	priv, err := crypto.GenerateKeyPair()
	if err != nil {
		return err
	}

	server.privateKey = &priv
	server.publicKey = server.privateKey.PublicKey

	return nil
}

func (server *GameServer) Init() bool {
	var status bool

	logErr := server.initLogger()
	if logErr != nil {
		slog.Error("Failed to open log file for saving")
		return status
	}

	cryptoErr := server.initCrypto()
	if cryptoErr != nil {

		if cryptoErr == arcaneErrors.ErrKeyGenerationFailed {
			slog.Error("Failed to generate RSA keys for server")
		} else if cryptoErr == arcaneErrors.ErrKeysNotValid {
			slog.Error("Failed to validate the generated keys for the server")
		}

		return status
	}

	server.URI = "127.0.0.1:" + viper.GetString("port")
	server.MaxConnections = 8

	status = true
	return status
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

func (server *GameServer) CloseConnection(conn net.Conn) {
	conn.Close()
	server.ConnectionCount -= 1
}

func (server *GameServer) Read(conn net.Conn) (string, error) {
	var ret string

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			slog.Info("Client has disconnected (EOF)", "client", conn.RemoteAddr().String())
		} else {
			slog.Error("Failed to read buffer from client", "err", err.Error(), "client", conn.RemoteAddr().String())
		}
		return ret, err
	}

	ret = string(buffer[:n])

	slog.Debug("Message from client", "msg", ret)
	return ret, nil
}

func (server *GameServer) ReadEncrypted(conn net.Conn) (string, error) {

	cipherText, err := server.Read(conn)
	if err != nil {
		return "", err
	}

	plainText, err := crypto.DecryptMessage(cipherText, server.privateKey)
	if err != nil {
		if err == arcaneErrors.ErrBase64DecodeFailed {
			slog.Error("Failed to decrypt base64 encoded cipher text")
		} else if err == arcaneErrors.ErrDecryptionFailed {
			slog.Error("Failed to decrypt cipher text provided by the client")
		}

		return "", err
	}

	return plainText, nil
}

func (server *GameServer) Write(message string, conn net.Conn) error {
	buffer := []byte(message)

	_, err := conn.Write(buffer)
	if err != nil {
		slog.Error("Failed to write buffer to client", "err", err.Error(), "client", conn.RemoteAddr().String())
		return err
	}

	return nil
}

func (server *GameServer) HandleClient(conn net.Conn) {
	if !server.NegotiateKeys(conn) {
		slog.Error("Key negotiation failed for client. Closing connection...")
		server.CloseConnection(conn)
		return
	}

	slog.Info("Key negotiation success for client", "client", conn.RemoteAddr().String())
	slog.Info("Client OK. Waiting for JOIN request", "client", conn.RemoteAddr().String())
	for {
		plainText, err := server.ReadEncrypted(conn)
		if err != nil {
			server.CloseConnection(conn)
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
	wErr := server.Write(keyResp, conn)
	if wErr != nil {
		return result
	}

	result = true
	return result
}

func (server *GameServer) Start() {
	initErr := server.Init()
	if !initErr {
		panic(initErr)
	}

	server.Listen()
	server.AcceptConnections()
}

func (server *GameServer) Stop() {
	sock := *server.Listener
	sock.Close()
}
