package profile

import profileSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/profile"

type ProfileHandler struct {
	ProfileService *profileSvc.ProfileService
}

func NewProfileHandler(profileService *profileSvc.ProfileService) *ProfileHandler {
	return &ProfileHandler{ProfileService: profileService}
}
