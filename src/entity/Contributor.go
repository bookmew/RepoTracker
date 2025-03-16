package entity

type Contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}
