package pair

type Response struct {
	OK          bool   `json:"ok"`
	BoundOrigin string `json:"boundOrigin"`
}
