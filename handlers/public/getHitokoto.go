package public

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm/utils"
	"hitokoto-go/consts"
	"hitokoto-go/global"
	"hitokoto-go/types"
	htutils "hitokoto-go/utils"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GetHitokotoQueryParams struct {
	categories []*types.MetaCategory // Will be converted in parseQueryParams()
	encode     string
	charset    string // Only valid when encode is text
	callback   string // Only valid when encode is set to js
	selector   string // Only valid when encode is set to js
	minLength  uint
	maxLength  uint

	// hitokoto-go specific
	orderBy string // "default", "visit", "like"
	desc    bool
}

func HandlerGetHitokoto(ctx *gin.Context) {

	// Parse query params
	q := parseQueryParams(ctx)

	// Rand category
	targetCategory := selectCategory(q)

	// Rand hitokoto
	hitokoto := selectHitokoto(q, targetCategory)

	// Encode hitokoto
	switch q.encode {
	case "text":
		// Auto mime type
		htRes := encodeHitokotoToString(q, hitokoto)
		// Check charset
		if q.charset == "gbk" { // I don't know why this is provided, but using GBK is NOT a good idea
			ctx.Header("Content-Type", "text/plain; charset=gbk")
			htReader := transform.NewReader(bytes.NewReader([]byte(htRes)), simplifiedchinese.GBK.NewEncoder())
			htGBK, err := ioutil.ReadAll(htReader)
			if err != nil {
				// I just don't know how to handle this error
				log.Println("[ERROR] Failed to parse hitokoto ", htRes, " into GBK: ", err)
				ctx.String(http.StatusInternalServerError, "Failed to parse hitokoto into GBK")
			}
			ctx.String(http.StatusOK, string(htGBK))
		} else {
			ctx.String(http.StatusOK, htRes)
		}
	case "js":
		ctx.Header("Content-Type", "application/javascript; charset=utf-8")
		ctx.String(http.StatusOK, encodeHitokotoToString(q, hitokoto))
	case "json":
		// Auto mime type
		ctx.JSON(http.StatusOK, hitokoto)
	default: // Fallback to JSON
		// Auto mime type
		ctx.JSON(http.StatusOK, hitokoto)
	}

	// Save visit record
	go recordVisitCount(&types.VisitHitokotoRecord{
		Category: targetCategory.Key,
		ID:       hitokoto.ID,
	})

}

// RandGetOne : Get one random hitokoto in text format (for start-up msg)
func RandGetOne() string {
	q := &GetHitokotoQueryParams{
		categories: []*types.MetaCategory{}, // Just random
		encode:     "text",
		charset:    "utf-8",
		minLength:  1,
		maxLength:  0, // Default no limitation
		orderBy:    "default",
		desc:       false,
	}

	// Rand category
	targetCategory := selectCategory(q)

	// Rand hitokoto
	hitokoto := selectHitokoto(q, targetCategory)

	return encodeHitokotoToString(q, hitokoto)
}

func parseQueryParams(ctx *gin.Context) *GetHitokotoQueryParams {
	var q GetHitokotoQueryParams

	// Filter valid categories
	categories := ctx.QueryArray("c")
	for _, category := range global.Meta.Categories {
		if utils.Contains(categories, category.Key) {
			q.categories = append(q.categories, &category)
		}
	}

	// Parse base query params
	q.encode = ctx.DefaultQuery("encode", "json")
	q.charset = ctx.DefaultQuery("charset", "utf-8")
	q.callback = ctx.Query("callback")
	q.selector = ctx.DefaultQuery("selector", ".hitokoto")
	q.orderBy = strings.ToLower(ctx.DefaultQuery("order_by", "default"))
	q.desc = strings.Contains(strings.ToLower(ctx.DefaultQuery("desc", "true")), "t") // t / true / T / True

	// Parse length limits
	minLen, err := strconv.ParseInt(ctx.Query("min_length"), 10, 64)
	if err != nil {
		minLen = 1
	}
	maxLen, err := strconv.ParseInt(ctx.Query("max_length"), 10, 64)
	if err != nil {
		maxLen = 0 // Invalid limitation
	}
	q.minLength = uint(minLen)
	q.maxLength = uint(maxLen)

	return &q

}

func selectCategory(q *GetHitokotoQueryParams) *types.MetaCategory {
	// Select category, using weighted random
	var targetCategory *types.MetaCategory
	if len(q.categories) == 0 {
		// Just random
		tcid := uint(rand.Intn(int(global.Meta.AllCount)))
		for _, c := range global.Meta.Categories {
			if tcid <= c.Counts {
				targetCategory = &c
				break
			}
			tcid -= c.Counts
		}
	} else {
		// Rand from specified categories
		var ac uint // AllCounts
		for _, c := range q.categories {
			ac += c.Counts
		}
		tcid := uint(rand.Intn(int(ac)))
		for _, c := range q.categories {
			if tcid <= c.Counts {
				targetCategory = c
				break
			}
			tcid -= c.Counts
		}
	}
	return targetCategory
}

func selectHitokoto(q *GetHitokotoQueryParams, targetCategory *types.MetaCategory) *types.Sentence {
	// Parse limitation
	limit := &htutils.SelectLimitation{
		Category:    targetCategory.Key,
		MinLen:      q.minLength,
		MaxLen:      q.maxLength,
		OrderBy:     q.orderBy,
		IsDescOrder: q.desc,
	}

	// Prefer Cache
	hitokoto := htutils.SelectHitokotoFromCache(limit)
	if hitokoto == nil {
		// Fallback to DB
		hitokoto = htutils.SelectHitokotoFromDB(limit)
	}

	// If DB fails, nothing to do, so no need to check for nil
	return hitokoto

}

func encodeHitokotoToString(q *GetHitokotoQueryParams, hitokoto *types.Sentence) string {
	// Encode hitokoto
	var encoded string
	switch q.encode {
	case "text":
		encoded = hitokoto.Hitokoto
	case "js":
		if q.callback != "" {
			encoded = jsCallbackEncode(q.callback, hitokoto.Hitokoto, q.selector)
		} else {
			encoded = jsEncode(hitokoto.Hitokoto, q.selector)
		}
	default:
		// ???
		encoded = hitokoto.Hitokoto
	}
	return encoded
}

func jsEncode(hitokoto string, selector string) string {

	return fmt.Sprintf(
		`(function hitokoto(){var hitokoto="%s";var dom=document.querySelector('%s');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;})()`,
		hitokoto,
		selector,
	)
}

func jsCallbackEncode(callback string, hitokoto string, selector string) string {

	return fmt.Sprintf(
		`;%s("%s");`,
		callback,
		strings.ReplaceAll(
			jsEncode(hitokoto, selector),
			`"`,
			`\"`,
		),
	)
}

func recordVisitCount(r *types.VisitHitokotoRecord) {

	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	json := jsoniter.ConfigCompatibleWithStandardLibrary // For better performance

	// Add to redis cache
	recordBytes, err := json.Marshal(r)
	if err != nil {
		log.Println("[ERROR] Failed to marshal visit record:", err)
		return
	}
	global.Redis.LPush(ctx, consts.VisitListName, recordBytes)

}
