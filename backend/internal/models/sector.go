package models

type Sector struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Symbols []string `json:"symbols"`
}
