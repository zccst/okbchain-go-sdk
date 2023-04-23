module github.com/okx/okbchain-go-sdk

go 1.16

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/ethereum/go-ethereum v1.11.6
	github.com/golang/mock v1.6.0
	github.com/okx/okbchain v0.1.1
	github.com/stretchr/testify v1.8.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/go-amino => github.com/okex/go-amino v0.15.1-okc4
)
