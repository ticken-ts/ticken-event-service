package dto

type Validator struct {
	ValidatorID    string `json:"validator_id"`
	CreatedBy      string `json:"created_by"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
}
