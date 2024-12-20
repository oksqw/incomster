package incomster

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target ./backend/api/oas --package oas --clean ./backend/api/openapi.yaml
//go:generate sqlboiler psql --config ./backend/store/postgres/generate/sqlboiler.toml
