package sessioncache

import (
	"fmt"
	"github.com/pmylund/go-cache"
	"time"
)

type GoCache struct {
	gocache *cache.Cache
}

func (g *GoCache) Get(key string) (interface{}, error) {

	value, found := g.gocache.Get(key)
	if !found {
		return nil, fmt.Errorf("Key %v not found in cache", key)
	}

	return value, nil

}

func (g *GoCache) Set(key string, value interface{}, cacheDuration int) {
	g.gocache.Set(key, value, -1)
}

// Initialise the gocache object and set up the cache
func (g *GoCache) New() {
	g.gocache = cache.New(5*time.Minute, 30*time.Second)
}
