package public

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm/utils"
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type GetHitokotoQueryParams struct {
	categories []string
	encode     string
	charset    string // Only valid when encode is text
	callback   string // Only valid when encode is set to js
	selector   string // Only valid when encode is set to js
	minLength  uint
	maxLength  uint
}

func GetHitokoto(ctx *gin.Context) {

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
		ctx.JSON(http.StatusOK, hitokoto.ToJSON(targetCategory.Key))
	default: // Fallback to JSON
		// Auto mime type
		ctx.JSON(http.StatusOK, hitokoto.ToJSON(targetCategory.Key))
	}

}

// RandGetOne : Get one random hitokoto in text format (for start-up msg)
func RandGetOne() string {
	q := &GetHitokotoQueryParams{
		categories: []string{}, // Just random
		encode:     "text",
		charset:    "utf-8",
		minLength:  0,
		maxLength:  60, // Default max length is 60
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
			q.categories = append(q.categories, category.Key)
		}
	}

	// Parse base query params
	q.encode = ctx.DefaultQuery("encode", "json")
	q.charset = ctx.DefaultQuery("charset", "utf-8")
	q.callback = ctx.Query("callback")
	q.selector = ctx.DefaultQuery("selector", ".hitokoto")

	// Parse length limits
	minLen, err := strconv.ParseInt(ctx.Query("min_length"), 10, 64)
	if minLen < 0 || err != nil {
		minLen = 0
	}
	maxLen, err := strconv.ParseInt(ctx.Query("max_length"), 10, 64)
	if err != nil {
		maxLen = 60 // Why 60? category L has min length of 31, so use 30 will cause record not found
	} else if maxLen < minLen {
		maxLen = minLen
	}
	q.minLength = uint(minLen)
	q.maxLength = uint(maxLen)

	return &q

}

func selectCategory(q *GetHitokotoQueryParams) *types.MetaCategory {
	// Select category
	var targetCategory types.MetaCategory
	if len(q.categories) == 0 {
		// Just random
		tcid := rand.Intn(len(global.Meta.Categories))
		targetCategory = global.Meta.Categories[tcid]
	} else {
		// Rand from specified categories
		tcid := rand.Intn(len(q.categories))
		tcKey := q.categories[tcid]
		// Find category
		for _, c := range global.Meta.Categories {
			if c.Key == tcKey {
				targetCategory = c
				break
			}
		}
	}
	return &targetCategory
}

func selectHitokoto(q *GetHitokotoQueryParams, targetCategory *types.MetaCategory) *models.Sentence {
	// Select hitokoto
	var ht models.Sentence
	global.DB.
		Scopes(
			models.SentenceTable(
				models.Sentence{Type: targetCategory.Key},
			),
		).
		Where("length >= ?", q.minLength).
		Where("length <= ?", q.maxLength).
		Order("RANDOM()").
		First(&ht)
	return &ht
}

func encodeHitokotoToString(q *GetHitokotoQueryParams, hitokoto *models.Sentence) string {
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
