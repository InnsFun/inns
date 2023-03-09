package inns_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "inns/testutil/keeper"
	"inns/testutil/nullify"
	"inns/x/simple"
	"inns/x/simple/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.InnsKeeper(t)
	inns.InitGenesis(ctx, *k, genesisState)
	got := inns.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
