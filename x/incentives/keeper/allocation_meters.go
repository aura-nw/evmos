// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/evmos/v15/x/incentives/types"
)

// GetAllAllocationMeters - get all registered AllocationMeters
func (k Keeper) GetAllAllocationMeters(ctx sdk.Context) []math.LegacyDecCoin {
	allocationMeters := []math.LegacyDecCoin{}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixAllocationMeter)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Key()[1:])
		allocationMeter, found := k.GetAllocationMeter(ctx, denom)
		if !found {
			continue
		}

		allocationMeters = append(allocationMeters, allocationMeter)
	}

	return allocationMeters
}

// GetAllocationMeter - get registered allocationMeter from the identifier
func (k Keeper) GetAllocationMeter(
	ctx sdk.Context,
	denom string,
) (math.LegacyDecCoin, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAllocationMeter)

	bz := store.Get([]byte(denom))
	if bz == nil {
		return math.LegacyDecCoin{
			Denom:  denom,
			Amount: math.LegacyZeroDec(),
		}, false
	}

	var amount math.LegacyDec
	err := amount.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal amount value %v", err))
	}
	return math.LegacyDecCoin{
		Denom:  denom,
		Amount: amount,
	}, true
}

// SetAllocationMeter stores an allocationMeter
func (k Keeper) SetAllocationMeter(ctx sdk.Context, am math.LegacyDecCoin) {
	bz, err := am.Amount.Marshal()
	if err != nil {
		panic(fmt.Errorf("unable to marshal amount value %v", err))
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAllocationMeter)
	key := []byte(am.Denom)

	// Remove zero allocation meters
	if am.IsZero() {
		store.Delete(key)
	} else {
		store.Set(key, bz)
	}
}
