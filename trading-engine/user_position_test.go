package trading_engine

import "testing"

func TestPosition(t *testing.T) {
	Deposit("1", "1", 100)
	Frozen("1", "1", 100)
	WithDraw("1", "1", 100, true)
	t.Log(UserAssetMap)
}
