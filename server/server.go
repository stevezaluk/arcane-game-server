package server

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stevezaluk/arcane-game-server/crypto"
	arcaneErrors "github.com/stevezaluk/arcane-game-server/errors"
	"io"
	"log/slog"
	"net"
	"strings"
)

type GameServer struct {
	URI             string
	Listener        *net.Listener
	ConnectionCount int
	IsClosed        bool

	ServerKeyPair crypto.KeyPair
}

func (server *GameServer) initCrypto() error {
	slog.Info("Generating RSA-4096 key pair...")
	keyPair, err := crypto.New()
	if err != nil {
		return err
	}

	server.ServerKeyPair = keyPair
	slog.Info("Key Pair", "key", server.ServerKeyPair.PublicKeyChecksum)

	return nil
}

func (server *GameServer) Init() bool {
	var status bool

	cryptoErr := server.initCrypto()
	if cryptoErr != nil {

		if errors.Is(cryptoErr, arcaneErrors.ErrKeyGenerationFailed) {
			slog.Error("Failed to generate RSA keys for server")
		} else if errors.Is(cryptoErr, arcaneErrors.ErrKeysNotValid) {
			slog.Error("Failed to validate the generated keys for the server")
		}

		return status
	}

	server.URI = "127.0.0.1:" + viper.GetString("port")

	status = true
	return status
}

func (server *GameServer) Listen() error {
	slog.Info("Starting server...")
	listen, err := net.Listen("tcp", server.URI)
	if err != nil {
		return arcaneErrors.ErrServerStartFailed
	}

	slog.Info("Server listening for connections at", "uri", server.URI)
	server.Listener = &listen

	return nil
}

func (server *GameServer) WaitForConnections() {
	for {
		if server.ConnectionCount == viper.GetInt("server.max_connections") {
			server.IsClosed = true
			slog.Warn("Server reached max connection count. No new connections are accepted")
			break
		}

		if !server.IsClosed {
			conn, err := server.AcceptConnection()
			if err != nil {
				continue
			}

			go server.HandleClient(conn)
		}
	}
}

func (server *GameServer) AcceptConnection() (net.Conn, error) {
	sock := *server.Listener

	conn, err := sock.Accept()
	if err != nil {
		slog.Error("Failed to accept client connection", "client", conn.RemoteAddr().String())
		conn.Close()
		return nil, arcaneErrors.ErrAcceptConnectionFailed
	}

	slog.Info("Client has connected", "client", conn.RemoteAddr().String())
	server.ConnectionCount += 1
	return conn, nil
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

	plainText, err := crypto.DecryptMessage(cipherText, server.ServerKeyPair.PrivateKey)
	if err != nil {
		if errors.Is(err, arcaneErrors.ErrBase64DecodeFailed) {
			slog.Error("Failed to decrypt base64 encoded cipher text")
		} else if errors.Is(err, arcaneErrors.ErrDecryptionFailed) {
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
	slog.Info("Starting key negotiation with client", "client", conn.RemoteAddr().String())

	err := server.NegotiateServerKey(conn)
	if err != nil {
		if errors.Is(err, arcaneErrors.ErrReadBufferFailed) {
			slog.Error("Failed to read buffer while waiting for connect response")
		} else if errors.Is(err, arcaneErrors.ErrWriteBufferFailed) {
			slog.Error("Failed to write public key to client")
		} else if errors.Is(err, arcaneErrors.ErrInvalidConnectResponse) {
			slog.Error("CONNECT Request not formatted properly", "client", conn.RemoteAddr().String())
		}

		slog.Error("Key negotiation failed for client. Closing connection...")
		server.CloseConnection(conn)
		return
	}

	slog.Info("Waiting for public key acknowledgement from client")

	err = server.ValidateServerKey(conn)
	if err != nil {
		if errors.Is(err, arcaneErrors.ErrReadBufferFailed) {
			slog.Error("Failed to read buffer while waiting for key acknowledgement")
		} else if errors.Is(err, arcaneErrors.ErrServerClientKeyMismatch) {
			slog.Error("Client responded with invalid public key checksum. Server and client key mismatch")
		}

		slog.Error("Server side key validation failed. Closing connection...")
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

func (server *GameServer) NegotiateServerKey(conn net.Conn) error {
	buffer, err := server.Read(conn)
	if err != nil {
		return arcaneErrors.ErrReadBufferFailed
	}

	if buffer != "CONNECT" {
		return arcaneErrors.ErrInvalidConnectResponse
	}

	slog.Info("CONNECT Response acknowledeged. PEM encoding public key...", "client", conn.RemoteAddr().String())

	keyResp := "PUBKEY:" + server.ServerKeyPair.PublicKeyPem

	slog.Info("Sending public key...")
	wErr := server.Write(keyResp, conn)
	if wErr != nil {
		return arcaneErrors.ErrWriteBufferFailed
	}

	return nil
}

func (server *GameServer) ValidateServerKey(conn net.Conn) error {
	buffer, err := server.Read(conn)
	if err != nil {
		return arcaneErrors.ErrReadBufferFailed
	}

	if !strings.HasPrefix(buffer, "PUBKEY:ACK:") {
		return arcaneErrors.ErrInvalidKeyAcknowledgement
	}

	clientChecksum := strings.Split(buffer, ":")[2]

	if clientChecksum != server.ServerKeyPair.PublicKeyChecksum {
		return arcaneErrors.ErrServerClientKeyMismatch
	}

	return nil
}

func (server *GameServer) Start() {
	initErr := server.Init()
	if !initErr {
		slog.Error("Server initialization has failed")
		panic(initErr)
	}

	listenErr := server.Listen()
	if listenErr != nil {
		slog.Error("Failed to start listening for connections", "err", listenErr.Error())
		panic(listenErr)
	}

	server.WaitForConnections()
}

func (server *GameServer) Stop() {
	sock := *server.Listener
	err := sock.Close()
	if err != nil {
		panic(err)
	}
}
