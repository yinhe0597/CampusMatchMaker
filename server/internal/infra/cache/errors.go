package cache

import "errors"

// ErrCacheMiss 缓存未命中（正常情况，调用方应回源查询）
var ErrCacheMiss = errors.New("cache: miss")
