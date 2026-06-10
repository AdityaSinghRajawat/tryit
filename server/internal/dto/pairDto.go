package dto

type PairRequest struct {
	Token string `json:"token"`
}

type PairResponse struct {
	OK          bool   `json:"ok"`
	BoundOrigin string `json:"boundOrigin"`
}
