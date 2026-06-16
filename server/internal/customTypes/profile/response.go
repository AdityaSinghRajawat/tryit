package profile

type ListResponse struct {
	Profiles []SiteProfile `json:"profiles"`
}

type CreateResponse struct {
	Host string `json:"host"`
}
