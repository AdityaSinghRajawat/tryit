package config

type pairConsts struct {
	tokenPrefix       string
	keychainEntryName string // Phase 2 use
}

var pairI = &pairConsts{
	tokenPrefix:       "tk_",
	keychainEntryName: "tryit-pair-token",
}

func GetTokenPrefix() string       { return pairI.tokenPrefix }
func GetKeychainEntryName() string { return pairI.keychainEntryName }
