package services

type organizationManager struct {
}

func NewOrganizationManager() OrganizationManager {
	return new(organizationManager)
}
