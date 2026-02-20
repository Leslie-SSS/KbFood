package platform

import (
	"context"
	"time"

	"kbfood/internal/domain/entity"
)

// Client represents a platform API client
type Client interface {
	// Name returns the platform name
	Name() string

	// FetchProducts fetches products from the platform (active mode)
	FetchProducts(ctx context.Context, region string) ([]*entity.PlatformProductDTO, error)

	// ShouldFetch checks if fetching is allowed at this time
	ShouldFetch(now time.Time) bool
}

// PushClient represents a platform that pushes data (passive mode)
type PushClient interface {
	// ProcessPushData processes pushed data from the platform
	ProcessPushData(ctx context.Context, data []*PushData) (map[string][]*entity.PlatformProductDTO, error)
}

// PushData represents pushed data from DT platform
type PushData struct {
	Title      string
	Price      float64
	OriginalPrice float64
	Status     int
	CrawlTime  int64
	Region     string
}

// RegionConfig holds latitude/longitude for regions
type RegionConfig struct {
	Name      string
	CityName  string
	Latitude  float64
	Longitude float64
}

// Default regions
var Regions = map[string]RegionConfig{
	"广州": {
		Name:      "广州",
		CityName:  "广州市",
		Latitude:  22.937719345092773,
		Longitude: 113.38423919677734,
	},
	"深圳": {
		Name:      "深圳",
		CityName:  "深圳市",
		Latitude:  22.5431,
		Longitude: 114.0579,
	},
	"北京": {
		Name:      "北京",
		CityName:  "北京市",
		Latitude:  39.9042,
		Longitude: 116.4074,
	},
	"上海": {
		Name:      "上海",
		CityName:  "上海市",
		Latitude:  31.2304,
		Longitude: 121.4737,
	},
	"成都": {
		Name:      "成都",
		CityName:  "成都市",
		Latitude:  30.5728,
		Longitude: 104.0668,
	},
	"重庆": {
		Name:      "重庆",
		CityName:  "重庆市",
		Latitude:  29.4316,
		Longitude: 106.9123,
	},
	"杭州": {
		Name:      "杭州",
		CityName:  "杭州市",
		Latitude:  30.2741,
		Longitude: 120.1551,
	},
	"武汉": {
		Name:      "武汉",
		CityName:  "武汉市",
		Latitude:  30.5928,
		Longitude: 114.3055,
	},
	"西安": {
		Name:      "西安",
		CityName:  "西安市",
		Latitude:  34.3416,
		Longitude: 108.9398,
	},
	"南京": {
		Name:      "南京",
		CityName:  "南京市",
		Latitude:  32.0603,
		Longitude: 118.7969,
	},
	"佛山": {
		Name:      "佛山",
		CityName:  "佛山市",
		Latitude:  23.0219,
		Longitude: 113.1214,
	},
}

// GetRegionConfig returns region config, defaults to Guangzhou
func GetRegionConfig(region string) RegionConfig {
	if cfg, ok := Regions[region]; ok {
		return cfg
	}
	return Regions["广州"]
}

// ConvertToCityName converts region to city name
func ConvertToCityName(region string) string {
	cfg := GetRegionConfig(region)
	return cfg.CityName
}
