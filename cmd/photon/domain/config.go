package domain

type Config struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	EntryPoint  string `json:"entryPoint"`
	DixConfig   string `json:"dixConfig"`
	CoreVersion string `json:"coreVer"`
}
