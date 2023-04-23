package slashing

import (
	"github.com/okx/okbchain-go-sdk/exposed"
	"github.com/okx/okbchain-go-sdk/module/slashing/types"
	gosdktypes "github.com/okx/okbchain-go-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/x/slashing"
)

var _ gosdktypes.Module = (*slashingClient)(nil)

type slashingClient struct {
	gosdktypes.BaseClient
}

// RegisterCodec registers the msg type in slashing module
func (sc slashingClient) RegisterCodec(cdc *codec.Codec) {
	slashing.RegisterCodec(cdc)
}

// Name returns the module name
func (slashingClient) Name() string {
	return types.ModuleName
}

// NewSlashingClient creates a new instance of slashing client as implement
func NewSlashingClient(baseClient gosdktypes.BaseClient) exposed.Slashing {
	return slashingClient{baseClient}
}
