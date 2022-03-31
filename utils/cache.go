package utils

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"log"
	"time"
)

type SelectLimitation struct {
	Category    string
	MinLen      uint
	MaxLen      uint
	OrderBy     string // "default", "visit", "like"
	IsDescOrder bool   // By default (LPush -> LPop)
}

func encodeCacheKey(limit *SelectLimitation) string {
	return fmt.Sprintf("hitokoto:cache:%s:sentence_%s", limit.OrderBy, limit.Category)
}

// CacheCategory : Pre-random sentences to prevent the database from being overloaded
func CacheCategory(c string) error {
	if !global.Config.IsProdMode {
		log.Println("[Info] Create cache for category ", c)
	}

	modes := []string{"default", "visit", "like"}
	for _, mode := range modes {
		limit := &SelectLimitation{
			Category:    c,
			OrderBy:     mode,
			IsDescOrder: true,
			MinLen:      1,
			MaxLen:      0,
		}

		// Get random table from database
		var mss []models.Sentence // Model sentences
		randReq := global.DB.
			Scopes(
				models.SentenceTable(
					models.Sentence{Type: c},
				),
			)

		switch mode {
		case "visit":
			randReq = randReq.Order("RANDOM()*(visit_count+1)") // +1 to prevent 0
		case "like":
			randReq = randReq.Order("RANDOM()*(like_count+1)") // +1 to prevent 0
		default:
			randReq = randReq.Order("RANDOM()")
		}

		randReq.
			// Limit(consts.RandTableSize). // Remove limit so that we can cache all sentences
			Find(&mss)

		// Transform to type
		var ss []*types.Sentence // Sentences
		for _, ms := range mss {
			ss = append(ss, ms.ToType(c))
		}

		if err := saveHitokotosIntoCache(limit, ss); err != nil {
			log.Println("[Error] Failed to cache category ", limit.Category, " on ", limit.OrderBy, " mode with error: ", err)
			return err
		}
	}

	return nil
}

func saveHitokotosIntoCache(limit *SelectLimitation, ss []*types.Sentence) error {

	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	json := jsoniter.ConfigCompatibleWithStandardLibrary // For better performance

	// Clear old cache
	global.Redis.LTrim(ctx, encodeCacheKey(limit), 1, 0)

	var cacheSentenceBytes []interface{} // ???
	for _, s := range ss {
		sentenceBytes, err := json.Marshal(&s)
		if err != nil {
			return err
		}
		cacheSentenceBytes = append(cacheSentenceBytes, sentenceBytes)
	}
	if len(cacheSentenceBytes) == 0 {
		// Nothing to cache
		log.Println("[WARN] Nothing to cache!")
		return nil
	}
	if n, err := global.Redis.LPush(ctx, encodeCacheKey(limit), cacheSentenceBytes...).Result(); err != nil {
		log.Println("[ERROR] Failed to save hitokoto into cache with error: ", err)
		return err
	} else {
		log.Println("[Info] Successfully saved ", n, " hitokotos of type ", limit.Category, " into cache")
	}

	return nil
}

func SelectHitokotoFromCache(limit *SelectLimitation) *types.Sentence {

	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	json := jsoniter.ConfigCompatibleWithStandardLibrary // For better performance

	// Check cache
	sCachedCount, err := global.Redis.LLen(ctx, encodeCacheKey(limit)).Result()
	if sCachedCount == 0 || err != nil {
		log.Println("[WARN] No hitokoto found in cache for category ", limit.Category)
		return nil
	}

	isMaxLengthLimitInvalid := limit.MaxLen < limit.MinLen

	var selectedSentence *types.Sentence
	pop := global.Redis.LPop
	if !limit.IsDescOrder {
		pop = global.Redis.RPop
	}
	push := global.Redis.RPush
	if !limit.IsDescOrder {
		push = global.Redis.LPush
	}

	for i := int64(0); i < sCachedCount; i++ {
		var s *types.Sentence
		cachedBytes, err := pop(ctx, encodeCacheKey(limit)).Bytes()
		if err != nil {
			log.Println("[ERROR] Fail to get cached category ", limit.Category, " index ", i, " with error: ", err)
			return nil
		}
		if err = json.Unmarshal(cachedBytes, &s); err != nil {
			log.Println("[ERROR] Fail to parse cached category ", limit.Category, " index ", i, " with error: ", err)
			return nil
		}

		// Push back to the end of the list
		if err = push(ctx, encodeCacheKey(limit), cachedBytes).Err(); err != nil {
			log.Println("[ERROR] Fail to push back cached category ", limit.Category, " with error: ", err)
		}

		// Check if the sentence is valid
		if s.Length >= limit.MinLen && // Check min length
			(isMaxLengthLimitInvalid || // If max length limitation is valid
				s.Length <= limit.MaxLen) { // Check max length
			// The one
			// Mark it as selected
			selectedSentence = s

			// Found!
			break
		}

	}
	return selectedSentence
}
