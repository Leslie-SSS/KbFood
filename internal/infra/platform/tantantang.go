package platform

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"kbfood/internal/config"
	"kbfood/internal/domain/entity"
)

// TanTanTangClient implements the TanTanTang platform client
type TanTanTangClient struct {
	cfg    *config.TanTanTangConfig
	client *resty.Client
}

// NewTanTanTangClient creates a new TanTanTang client
func NewTanTanTangClient(cfg *config.TanTanTangConfig) *TanTanTangClient {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetHeader("User-Agent",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080712) UnifiedPCMacWechat(0xf2641210) XWEB/16816").
		SetHeader("Accept", "*/*").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetHeader("Content-Type", "application/x-www-form-urlencoded")

	return &TanTanTangClient{
		cfg:    cfg,
		client: client,
	}
}

// Name returns the platform name
func (c *TanTanTangClient) Name() string {
	return "探探糖"
}

// ShouldFetch checks if fetching is allowed at this time
func (c *TanTanTangClient) ShouldFetch(now time.Time) bool {
	hour, min, _ := now.Clock()
	// Avoid maintenance window: 23:55 - 00:10
	if hour == 23 && min >= 55 {
		return false
	}
	if hour == 0 && min < 10 {
		return false
	}
	return true
}

// FetchProducts fetches products from TanTanTang
func (c *TanTanTangClient) FetchProducts(ctx context.Context, region string) ([]*entity.PlatformProductDTO, error) {
	if !c.ShouldFetch(time.Now()) {
		log.Info().Msg("Skipping fetch due to maintenance window")
		return nil, nil
	}

	var allProducts []*entity.PlatformProductDTO
	page := 1
	maxPages := 100
	maxRetries := 3

	for page <= maxPages {
		var products []*entity.PlatformProductDTO
		var hasMore bool
		var err error

		// Retry logic for transient failures
		for retry := 0; retry < maxRetries; retry++ {
			products, hasMore, err = c.fetchPage(ctx, region, page)
			if err == nil {
				break
			}

			log.Warn().
				Err(err).
				Int("page", page).
				Int("retry", retry+1).
				Int("maxRetries", maxRetries).
				Msg("Fetch page failed, retrying...")

			if retry < maxRetries-1 {
				time.Sleep(time.Duration(retry+1) * time.Second) // Exponential backoff
			}
		}

		if err != nil {
			// If all retries failed, log and continue to next page or return partial results
			log.Error().
				Err(err).
				Int("page", page).
				Int("productsSoFar", len(allProducts)).
				Msg("All retries exhausted for page")
			// Return partial results instead of failing completely
			if len(allProducts) > 0 {
				return allProducts, nil
			}
			return nil, fmt.Errorf("fetch page %d: %w", page, err)
		}

		allProducts = append(allProducts, products...)

		log.Info().
			Int("page", page).
			Int("pageProducts", len(products)).
			Int("totalProducts", len(allProducts)).
			Bool("hasMore", hasMore).
			Msg("Fetched page")

		if !hasMore {
			break
		}

		page++
		time.Sleep(1 * time.Second) // Rate limiting
	}

	log.Info().
		Int("totalProducts", len(allProducts)).
		Int("totalPages", page).
		Msg("Fetch completed")

	return allProducts, nil
}

// fetchPage fetches a single page of products
func (c *TanTanTangClient) fetchPage(ctx context.Context, region string, page int) ([]*entity.PlatformProductDTO, bool, error) {
	regionCfg := GetRegionConfig(region)
	cityName := regionCfg.CityName

	rqToken := c.generateSign(page, cityName)
	cityEncoded := url.QueryEscape(cityName)

	// Try multiple category IDs to find products
	// Common categories: 0=all, 14=food, etc.
	body := fmt.Sprintf(
		"cate_id=&cate2_id=0&area=&street=&page=%d&count=10&lon=%.14f&lat=%.15f&city=%s&rqtoken=%s",
		page, regionCfg.Longitude, regionCfg.Latitude, cityEncoded, rqToken,
	)

	// Safe token preview for logging
	tokenPreview := c.cfg.Token
	if len(tokenPreview) > 10 {
		tokenPreview = tokenPreview[:10] + "..."
	}

	log.Debug().
		Str("region", region).
		Str("city", cityName).
		Int("page", page).
		Str("token", tokenPreview).
		Str("rqToken", rqToken).
		Msg("API request")

	var resp APIResponse
	apiResp, err := c.client.R().
		SetContext(ctx).
		SetHeader("token", "").
		SetHeader("Host", "ttt.bjlxkjyxgs.cn").
		SetHeader("xweb_xhr", "1").
		SetHeader("sec-fetch-site", "cross-site").
		SetHeader("sec-fetch-mode", "cors").
		SetHeader("sec-fetch-dest", "empty").
		SetHeader("referer", "https://servicewechat.com/wx454addfc6819a2ac/117/page-frame.html").
		SetHeader("accept-language", "zh-CN,zh;q=0.9").
		SetBody(body).
		SetResult(&resp).
		SetDebug(true).
		Post(c.cfg.BaseURL)

	if err != nil {
		log.Error().Err(err).Msg("HTTP request failed")
		return nil, false, fmt.Errorf("http request: %w", err)
	}

	// Log raw response body for debugging
	rawBody := string(apiResp.Body())
	if len(resp.Data.Data) == 0 {
		log.Debug().
			Str("rawResponse", rawBody).
			Msg("Empty API response - raw body")
	} else {
		log.Debug().
			Int("productCount", len(resp.Data.Data)).
			Msg("Got products")
	}

	log.Debug().
		Int("statusCode", apiResp.StatusCode()).
		Int("code", resp.Code).
		Str("msg", resp.Msg).
		Int("dataCount", len(resp.Data.Data)).
		Msg("API response")

	if resp.Code != 1 {
		// Detailed error logging for different error codes
		var errorType string
		switch resp.Code {
		case -1:
			errorType = "AUTH_ERROR"
		case 0:
			errorType = "GENERIC_ERROR"
		case 401:
			errorType = "UNAUTHORIZED"
		case 403:
			errorType = "FORBIDDEN"
		case 429:
			errorType = "RATE_LIMITED"
		case 500:
			errorType = "SERVER_ERROR"
		default:
			errorType = "UNKNOWN_ERROR"
		}

		log.Warn().
			Int("code", resp.Code).
			Str("msg", resp.Msg).
			Str("errorType", errorType).
			Str("rawResponse", rawBody).
			Str("token", tokenPreview).
			Str("rqToken", rqToken).
			Msg("API returned error")

		// Return more specific error message
		return nil, false, fmt.Errorf("api error [%s]: code=%d, msg=%s", errorType, resp.Code, resp.Msg)
	}

	products := c.parseProducts(resp.Data.Data)
	log.Debug().
		Int("productsParsed", len(products)).
		Msg("Parsed products")

	hasMore := len(products) > 0

	return products, hasMore, nil
}

// generateSign generates the MD5 signature
func (c *TanTanTangClient) generateSign(page int, city string) string {
	// Build param string WITHOUT URL encoding (matching original JS implementation)
	// Format: c=城市&p=页码 (sorted alphabetically)
	paramStr := fmt.Sprintf("c=%s&p=%d", city, page)

	dateStr := time.Now().Format("20060102")
	signStr := c.cfg.SecretKey + "#/api/shop/activity-" + paramStr + "~" + dateStr

	// Debug: Log signature generation details (mask secret key)
	log.Debug().
		Str("paramStr", paramStr).
		Str("dateStr", dateStr).
		Str("signStrPreview", "****#/api/shop/activity-"+paramStr+"~"+dateStr).
		Msg("Generating signature")

	hash := md5.Sum([]byte(signStr))
	signature := hex.EncodeToString(hash[:])

	log.Debug().Str("signature", signature).Msg("Generated signature")

	return signature
}

// parseProducts parses API response into products
func (c *TanTanTangClient) parseProducts(items []ProductItem) []*entity.PlatformProductDTO {
	products := make([]*entity.PlatformProductDTO, 0, len(items))

	for _, item := range items {
		// Parse string prices to float64
		originalPrice, _ := strconv.ParseFloat(item.YPrice, 64)
		currentPrice, _ := strconv.ParseFloat(item.Price, 64)

		products = append(products, &entity.PlatformProductDTO{
			ActivityID:         strconv.FormatInt(item.ActivityID, 10),
			Platform:           c.Name(),
			Region:             "",
			Title:              item.Title,
			ShopName:           item.ShopName,
			OriginalPrice:      originalPrice,
			CurrentPrice:       currentPrice,
			SalesStatus:        item.SyStore,
			ActivityCreateTime: time.Unix(item.CreateTime, 0),
		})
	}

	return products
}

// API types
type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Data []ProductItem `json:"data"`
	} `json:"data"`
}

type ProductItem struct {
	ActivityID int64  `json:"activity_id"`
	Title      string `json:"title"`
	ShopName   string `json:"shop_name"`
	YPrice     string `json:"y_price"`
	Price      string `json:"price"`
	SyStore    int    `json:"sy_store"`
	CreateTime int64  `json:"createtime"`
}
