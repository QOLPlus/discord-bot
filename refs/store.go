package refs

import (
	coreAsset "github.com/QOLPlus/core/commands/asset"
	coreCrypto "github.com/QOLPlus/core/commands/cryptocurrency"
	"time"
)

type Store struct {
	Asset *AssetStore
	Ticker *TickerStore
}

type AssetStore struct {
	Timestamp int64
	Data      *[]AssetData
}

type TickerStore struct {
	Timestamp int64
	Data      *[]coreCrypto.MarketMaster
}

type AssetData struct {
	Master *coreAsset.StockMaster
	Assets *[]coreAsset.Asset
}

func (as *AssetStore) Reload(offset int64) {
	now := time.Now().Unix()
	if (now - as.Timestamp) < offset {
		return
	}

	masters, err := coreAsset.FetchStockMasters()
	if err != nil {
		return
	}

	var data []AssetData
	for _, master := range *masters {
		var assetData AssetData
		assets, err := master.FetchAssets()
		if err != nil {
			continue
		}
		assetData.Master = &master
		assetData.Assets = assets
		data = append(data, assetData)
	}

	if len(data) > 0 {
		as.Data = &data
		as.Timestamp = now
	}
}
func (ts *TickerStore) Reload(offset int64) {
	now := time.Now().Unix()
	if (now - ts.Timestamp) < offset {
		return
	}

	masters, err := coreCrypto.FetchMarketMasters()
	if err != nil {
		return
	}

	if len(*masters) > 0 {
		ts.Data = masters
		ts.Timestamp = now
	}
}