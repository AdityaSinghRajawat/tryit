package config

type cacheConsts struct {
	keyPrefix string
}

var cacheI = &cacheConsts{
	keyPrefix: "tryit:parse:",
}

func GetCacheKeyPrefix() string { return cacheI.keyPrefix }
