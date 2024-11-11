package errors

import "errors"

// Server Errors
var ErrServerStartFailed = errors.New("server: Failed to start listening for connections")
var ErrReadBufferFailed = errors.New("server: Failed to read buffer from client")
var ErrWriteBufferFailed = errors.New("server: Failed to write buffer to client")
var ErrAcceptConnectionFailed = errors.New("server: Failed to accept client connection")
var ErrCloseConnectionFailed = errors.New("server: Failed to close connect from client")
var ErrMaxConnectionsReached = errors.New("server: Cannot accept client connection. Server is full")

// Key Negotiation
var ErrInvalidConnectResponse = errors.New("keyNeogtiation: Did not receive expected CONNECT request from client during key negotiation")

// Decryption
var ErrKeysNotValid = errors.New("crypto: The keys generated for the server are not valid")
var ErrBase64DecodeFailed = errors.New("crypto: Failed to base64 decode encrypted message sent from client")
var ErrDecryptionFailed = errors.New("crypto: Failed to decrypt message from client")
