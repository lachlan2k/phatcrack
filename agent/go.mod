module github.com/lachlan2k/phatcrack/agent

go 1.23.0

toolchain go1.23.2

replace github.com/lachlan2k/phatcrack/common v0.0.0 => ../common

replace gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.6.0

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/lachlan2k/phatcrack/common v0.0.0
	github.com/nxadm/tail v1.4.11
	golang.org/x/sys v0.33.0
)

require (
	github.com/NHAS/webauthn v0.0.0-20240606085832-ea3172ef4dfa // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fxamacker/cbor/v2 v2.6.0 // indirect
	github.com/go-webauthn/x v0.1.10 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
