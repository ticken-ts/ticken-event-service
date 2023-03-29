package config

type DevUser struct {
	Email     string `mapstructure:"email"`
	UserID    string `mapstructure:"user_id"`
	Username  string `mapstructure:"username"`
	Firstname string `mapstructure:"firstname"`
	Lastname  string `mapstructure:"lastname"`
}

type MockInfo struct {
	DisablePVTBCMock bool `mapstructure:"disable_pvtbc_mock"`
	DisableBusMock   bool `mapstructure:"disable_bus_mock"`
	DisableAuthMock  bool `mapstructure:"disable_auth_mock"`
}

type Orgs struct {
	TickenOrgName string `mapstructure:"ticken_org_name"`
	TotalFakeOrgs int    `mapstructure:"total_fake_orgs"`
}

type Events struct {
	EventID          string         `mapstructure:"event_id"`
	EventName        string         `mapstructure:"event_name"`
	EventDescription string         `mapstructure:"event_description"`
	EventDate        string         `mapstructure:"event_date"`
	EventPosterUri   string         `mapstructure:"event_poster_uri"`
	EventSections    []EventSection `mapstructure:"event_sections"`
}

type EventSection struct {
	SectionName     string  `mapstructure:"section_name"`
	SectionPrice    float64 `mapstructure:"section_price"`
	SectionQuantity int     `mapstructure:"section_quantity"`
}

type DevConfig struct {
	User          DevUser  `mapstructure:"user"`
	Mock          MockInfo `mapstructure:"mock"`
	Orgs          Orgs     `mapstructure:"orgs"`
	Events        Events   `mapstructure:"events"`
	JWTPublicKey  string   `mapstructure:"jwt_public_key"`
	JWTPrivateKey string   `mapstructure:"jwt_private_key"`
}
