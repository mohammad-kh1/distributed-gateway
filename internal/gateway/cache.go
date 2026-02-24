package gateway

import (
	"time"

	"github.com/mohammad-kh1/distributed-gateway/api/proto"
	"github.com/patrickmn/go-cache"
)

type AuthCache struct {
	store *cache.Cache
}

func NewAuthCache() *AuthCache {
	// cache for 1 minute and delete in every 2 minute
	return &AuthCache{
		store: cache.New(1*time.Minute, 2*time.Minute),
	}
}

func (c *AuthCache) Get(token string) (*proto.VerifyResponse, bool) {
	val, found := c.store.Get(token)
	if !found {
		return nil, false
	}
	return val.(*proto.VerifyResponse), true
}

func (c *AuthCache) Set(token string, res *proto.VerifyResponse) {
	c.store.Set(token, res, cache.DefaultExpiration)
}
