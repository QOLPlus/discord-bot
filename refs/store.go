package refs

import (
	coreAsset "github.com/QOLPlus/core/commands/asset"
	"time"
)

type Store struct {
	Asset AssetStore
}

type AssetStore struct {
	Timestamp int64
	Data      []AssetData
}

type AssetData struct {
	Master coreAsset.StockMaster
	Assets []coreAsset.Asset
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
		assetData.Master = master
		assetData.Assets = *assets
		data = append(data, assetData)
	}

	if len(data) > 0 {
		as.Data = data
		as.Timestamp = now
	}
}
