package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"inns/x/mint/types"
)

// Simulation parameter constants
const (
	InitInflationRate   = "init_inflation_rate"
)

// GenInflationRateChange randomized InflationRateChange
func GenInitInflationRate(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var initInflationRate sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InitInflationRate, &initInflationRate, simState.Rand,
		func(r *rand.Rand) { initInflationRate = GenInitInflationRate(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	startTime := time.Now().AddDate(0, 0, 0)
	reductionFactor := sdk.NewDec(4).QuoInt64(5)
	initialAnnualProvisions := sdk.NewDec(1_000_000_000_000_000)
	params := types.NewParams(mintDenom,startTime, reductionFactor, initInflationRate, initialAnnualProvisions, blocksPerYear)

	mintGenesis := types.NewGenesisState(types.InitialMinter(), params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
