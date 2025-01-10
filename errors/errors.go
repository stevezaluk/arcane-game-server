package errors

import "errors"

// Logging Errors
var ErrLogFileFailed = errors.New("log: Failed to open log file")

// Server Errors
var ErrServerStartFailed = errors.New("server: Failed to start listening for connections")
var ErrReadBufferFailed = errors.New("server: Failed to read buffer from client")
var ErrWriteBufferFailed = errors.New("server: Failed to write buffer to client")
var ErrAcceptConnectionFailed = errors.New("server: Failed to accept client connection")
var ErrCloseConnectionFailed = errors.New("server: Failed to close connect from client")
var ErrMaxConnectionsReached = errors.New("server: Cannot accept client connection. Server is full")
var ErrParsePubKeyFailed = errors.New("server: Failed to parse public key")

// Key Negotiation
var ErrInvalidConnectResponse = errors.New("keyNegotiation: Did not receive expected CONNECT request from client during key negotiation")
var ErrInvalidKeyAcknowledgement = errors.New("keyNegotiation: Did not receive expected key acknowledgement from client or server during key negotiation")

// Decryption
var ErrKeyGenerationFailed = errors.New("crypto: Failed to generate keys for the server")
var ErrServerClientKeyMismatch = errors.New("crypto: Client key is different than server keys")
var ErrKeysNotValid = errors.New("crypto: The keys generated for the server are not valid")
var ErrBase64DecodeFailed = errors.New("crypto: Failed to base64 decode encrypted message sent from client")
var ErrDecryptionFailed = errors.New("crypto: Failed to decrypt message from client")

// Zone
var ErrZoneCannotBeShared = errors.New("zone: A zone cannot have an owner and be shared")
