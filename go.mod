module github.com/sigwinch28/miniserve

go 1.16

require (
	aead.dev/minisign v0.1.2
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/gorilla/handlers v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
)

replace aead.dev/minisign => ./deps/aead.dev/minisign
