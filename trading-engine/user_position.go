package trading_engine

import "errors"

type Asset struct {
	Available float64
	Frozen    float64
}

var UserAssetMap = map[string]map[string]Asset{}

type TransferType string

const AvailableToAvailable TransferType = "availableToAvailable"
const AvailableToFrozen TransferType = "availableToFrozen"
const FrozenToAvailable TransferType = "frozenToAvailable"

var TransferTypeError = errors.New("transfer type error")
var AmountError = errors.New("amount valid")

func TryTransfer(transferType TransferType, from, to, assetId string, amount float64, checkAmount bool) (bool, error) {
	if UserAssetMap[from] == nil {
		UserAssetMap[from] = make(map[string]Asset)
	}
	if UserAssetMap[to] == nil {
		UserAssetMap[to] = make(map[string]Asset)
	}
	switch transferType {
	case AvailableToFrozen:
		if checkAmount && UserAssetMap[from][assetId].Available < amount {
			return false, AmountError
		}
		asset := UserAssetMap[from][assetId]
		asset.Available = asset.Available - amount
		UserAssetMap[from][assetId] = asset

		toAsset := UserAssetMap[to][assetId]
		toAsset.Frozen = toAsset.Frozen + amount
		UserAssetMap[to][assetId] = toAsset
	case AvailableToAvailable:
		if checkAmount && UserAssetMap[from][assetId].Available < amount {
			return false, AmountError
		}
		fromAsset := UserAssetMap[from][assetId]
		fromAsset.Available = fromAsset.Available - amount
		UserAssetMap[from][assetId] = fromAsset

		toAsset := UserAssetMap[to][assetId]
		toAsset.Available = toAsset.Available + amount
		UserAssetMap[to][assetId] = toAsset
	case FrozenToAvailable:
		if checkAmount && UserAssetMap[from][assetId].Frozen < amount {
			return false, AmountError
		}
		fromAsset := UserAssetMap[from][assetId]
		fromAsset.Frozen = fromAsset.Frozen - amount
		UserAssetMap[from][assetId] = fromAsset

		toAsset := UserAssetMap[to][assetId]
		toAsset.Available = toAsset.Available + amount
		UserAssetMap[to][assetId] = toAsset
	default:
		return false, TransferTypeError
	}
	return true, nil
}

func Frozen(from, assetId string, amount float64) (bool, error) {
	return TryTransfer(AvailableToFrozen, from, from, assetId, amount, true)
}

func Deposit(to, assetId string, amount float64) {
	TryTransfer(AvailableToAvailable, "-1", to, assetId, amount, false)
}

func WithDraw(from, assetId string, amount float64, checkAmount bool) (bool, error) {
	return TryTransfer(FrozenToAvailable, from, "-1", assetId, amount, true)
}
