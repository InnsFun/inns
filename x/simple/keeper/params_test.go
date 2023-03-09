package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "inns/testutil/keeper"
	"inns/x/simple/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.InnsKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
