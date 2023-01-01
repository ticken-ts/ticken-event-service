package dto

type OrganizerDTO struct {
	OrganizerID string `json:"organizer_id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}
