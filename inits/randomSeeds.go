package inits

import (
	"math/rand"
	"time"
)

func RandomSeeds() error {
	rand.Seed(time.Now().UnixNano())
	return nil
}
