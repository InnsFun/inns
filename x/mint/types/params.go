package types

import (
	"errors"
	"fmt"
	"strings"
	time "time"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMintDenom           = []byte("MintDenom")
	KeyStartTime           = []byte("StartTime")
	KeyReductionFactor     = []byte("ReductionFactor")
	KeyInitInflationRate   = []byte("InitInflationRate")
	KeyInitialAnnualProvisions = []byte("InitialAnnualProvisions")
	KeyBlocksPerYear       = []byte("BlocksPerYear")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, startTime time.Time, reductionFactor, initInflationRate, initialAnnualProvisions sdk.Dec, blocksPerYear uint64,
) Params {
	return Params{
		MintDenom:           mintDenom,
		StartTime:           startTime,
		ReductionFactor:     reductionFactor,
		InitInflationRate: 	 initInflationRate,
		InitialAnnualProvisions: initialAnnualProvisions,
		BlocksPerYear:       blocksPerYear,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:           sdk.DefaultBondDenom,
		StartTime:           time.Now(),
		ReductionFactor:     sdk.NewDec(4).QuoInt64(5),   // 4/5
		InitInflationRate:   sdk.NewDecWithPrec(10, 2),   // 20%
		InitialAnnualProvisions: sdk.NewDec(1_000_000_000_000_000), // 1B
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateStartTime(p.StartTime); err != nil {
		return err
	}
	if err := validateReductionFactor(p.ReductionFactor); err != nil {
		return err
	}
	if err := validateInitInflationRate(p.InitInflationRate); err != nil {
		return err
	}
	if err := validateStartProvisions(p.InitialAnnualProvisions); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyStartTime, &p.StartTime, validateStartTime),
		paramtypes.NewParamSetPair(KeyReductionFactor, &p.ReductionFactor, validateReductionFactor),
		paramtypes.NewParamSetPair(KeyInitInflationRate, &p.InitInflationRate, validateInitInflationRate),
		paramtypes.NewParamSetPair(KeyInitialAnnualProvisions, &p.InitialAnnualProvisions, validateStartProvisions),
		paramtypes.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateStartTime(i interface{}) error {
	v, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsZero() {
		return fmt.Errorf("start time cannot be zero value: %s", v)
	}

	return nil
}

func validateReductionFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDec(1)) {
		return fmt.Errorf("reduction factor cannot be greater than 1")
	}

	if v.IsNegative() {
		return fmt.Errorf("reduction factor cannot be negative")
	}

	return nil
}

func validateInitInflationRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate change cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("inflation rate change too large: %s", v)
	}

	return nil
}

func validateStartProvisions(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("start provisions cannot be negative")
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}
