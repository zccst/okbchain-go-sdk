package tendermint

import (
	"github.com/okx/okbchain-go-sdk/exposed"
	"github.com/okx/okbchain-go-sdk/module/tendermint/types"
	gosdktypes "github.com/okx/okbchain-go-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
)

var _ gosdktypes.Module = (*tendermintClient)(nil)

type tendermintClient struct {
	gosdktypes.BaseClient
}

// nolint
func (tendermintClient) RegisterCodec(_ *codec.Codec) {}

// Name returns the module name
func (tendermintClient) Name() string {
	return types.ModuleName
}

// NewTendermintClient creates a new instance of tendermint client as implement
func NewTendermintClient(baseClient gosdktypes.BaseClient) exposed.Tendermint {
	return tendermintClient{baseClient}
}
