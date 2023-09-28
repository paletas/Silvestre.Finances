package exchanges

import "github.com/paletas/silvestre.finances/internal/pkg/assets"

type Exchange struct {
	ID        int64
	Name      string
	AssetType []assets.AssetType
}
