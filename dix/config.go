package dix

type Provider struct {
	From         string   `json:"from"`
	Factory      string   `json:"factory"`
	Dependencies []string `json:"dependencies"`
}

type Config struct {
	Module    string               `json:"module"`
	Pkg       string               `json:"pkg"`
	Imports   []string             `json:"imports"`
	Providers map[string]*Provider `json:"providers"`
}
