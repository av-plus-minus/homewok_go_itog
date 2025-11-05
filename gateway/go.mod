module gateway

go 1.25

require (
	github.com/av-plus-minus/homewok_go_05/ledger v0.0.0-20251028221037-3e75b7881f28
	github.com/go-chi/chi/v5 v5.2.3
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

replace github.com/av-plus-minus/homewok_go_05/ledger => ../ledger
