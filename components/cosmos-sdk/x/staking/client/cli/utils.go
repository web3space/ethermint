package cli

import (
	"errors"

	sdk "web3space/ethermint/components/cosmos-sdk/types"
	"web3space/ethermint/components/cosmos-sdk/x/staking/types"
)

func buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr string) (commission types.CommissionRates, err error) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, errors.New("must specify all validator commission parameters")
	}

	rate, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}

	maxRate, err := sdk.NewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}

	maxChangeRate, err := sdk.NewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}

	commission = types.NewCommissionRates(rate, maxRate, maxChangeRate)
	return commission, nil
}