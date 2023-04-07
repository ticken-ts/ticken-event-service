package fakes

import "fmt"

type SeedOrganizer struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (loader *Loader) seedOrganizer(toSeed []*SeedOrganizer) []error {
	var seedErrors = make([]error, 0)

	if loader.repoProvider.GetOrganizerRepository().Count() > 0 {
		return seedErrors
	}

	for _, organizer := range toSeed {
		_, err := loader.serviceProvider.GetOrganizerManager().RegisterOrganizer(
			organizer.Username,
			organizer.Username, // password same as username
			organizer.Email,
			organizer.Firstname,
			organizer.Lastname,
		)

		if err != nil {
			seedErrors = append(
				seedErrors,
				fmt.Errorf("failed to seed organizer with username %s: %s", organizer.Username, err.Error()),
			)
		}
	}

	return seedErrors
}
