package profile

type CreateRequest = SiteProfile

type CreateResponse struct {
	Host string `json:"host"`
}
