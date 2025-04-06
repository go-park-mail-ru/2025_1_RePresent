module retarget-authapp

go 1.24.0

require (
	github.com/go-playground/validator/v10 v10.26.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.7.3
	golang.org/x/crypto v0.36.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

require (
	gopkg.in/yaml.v3 v3.0.1
	pkg v0.0.0-00010101000000-000000000000
)

replace pkg => ../pkg
