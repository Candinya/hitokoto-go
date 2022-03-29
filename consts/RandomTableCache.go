package consts

import "time"

const (
	RandTableCacheExpire       = 60 * time.Minute
	RandTableSize              = 100
	RandTableRecreateThreshold = 5
)
