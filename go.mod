module openai-example

go 1.24.0

require github.com/tinfoilsh/encrypted-http-body-protocol v0.0.0

require (
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace github.com/tinfoilsh/encrypted-http-body-protocol => ../encrypted-http-body-protocol
