package config

type pairConsts struct {
	tokenPrefix       string
	keychainEntryName string
}

var pairI = &pairConsts{
	tokenPrefix:       "tk_",
	keychainEntryName: "tryit-pair-token",
}

func GetTokenPrefix() string       { return pairI.tokenPrefix }
func GetKeychainEntryName() string { return pairI.keychainEntryName }
