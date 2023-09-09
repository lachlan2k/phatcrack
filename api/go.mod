module github.com/lachlan2k/phatcrack/api

go 1.21

require (
	github.com/NHAS/webauthn v0.0.0-20230701002608-24fb1253febd
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/google/uuid v1.3.1
	github.com/gorilla/websocket v1.5.0
	github.com/labstack/echo/v4 v4.11.1
	github.com/lachlan2k/phatcrack/common v0.0.0
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.13.0
	gorm.io/datatypes v1.1.2-0.20230323024724-8e2e3c689dc8
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.4
)

require (
	github.com/fxamacker/cbor/v2 v2.5.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-webauthn/x v0.1.4 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gorm.io/driver/mysql v1.5.1 // indirect
)

replace github.com/lachlan2k/phatcrack/common v0.0.0 => ../common
