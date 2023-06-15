module github.com/axelarnetwork/utils

go 1.18

require (
	github.com/go-errors/errors v1.4.2
	github.com/matryer/moq v0.3.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
	github.com/tendermint/tendermint v0.34.2
	golang.org/x/exp v0.0.0-20221018221608-02f3b879a704
	golang.org/x/sync v0.1.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/tools v0.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tendermint/tendermint => github.com/cometbft/cometbft v0.34.27
