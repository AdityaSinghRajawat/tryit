package config

type cacheConsts struct {
	lruCapacity int
	diskSubdir  string
}

var cacheI = &cacheConsts{
	lruCapacity: 200,
	diskSubdir:  "/.tryit/cache",
}

func GetCacheLRUCapacity() int { return cacheI.lruCapacity }

// GetCacheDiskDir resolves the on-disk parse cache directory or "" when
// $HOME is unresolved.
func GetCacheDiskDir() string {
	if envConfigI == nil || envConfigI.homeDir == "" {
		return ""
	}
	return envConfigI.homeDir + cacheI.diskSubdir
}
