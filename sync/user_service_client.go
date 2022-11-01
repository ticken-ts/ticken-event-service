package sync

type UserServiceClient struct {
}

type UserMembership struct {
	MspID          string
	PeerEndpoint   string
	GatewayPeer    string
	Certificate    string
	PrivateKey     string
	TLSCertificate string
}

type UserInfo struct {
	Username       string
	UserID         string
	OrganizationID string
	IsAdmin        bool
}

func NewUserServiceClient() *UserServiceClient {
	return new(UserServiceClient)
}

const Certificate = "-----BEGIN CERTIFICATE-----\nMIICKzCCAdGgAwIBAgIRALLCYSWmrj8HhymZZjPnVqgwCgYIKoZIzj0EAwIwczEL\nMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG\ncmFuY2lzY28xGTAXBgNVBAoTEG9yZzIuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh\nLm9yZzIuZXhhbXBsZS5jb20wHhcNMjIxMTAxMDIwOTAwWhcNMzIxMDI5MDIwOTAw\nWjBsMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN\nU2FuIEZyYW5jaXNjbzEPMA0GA1UECxMGY2xpZW50MR8wHQYDVQQDDBZVc2VyMUBv\ncmcyLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE78fLFnws\nknPIe7zUXRG7/I0pofIlAwvhQXFiB8bF28nXfhiVXatD1H0JsZC7K37XwXEWE6fV\nO4piBgwx04Wc3KNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYD\nVR0jBCQwIoAgSbEwCE7s6OaDxuhk+ltHOo6eTOTRLubfmcpCdd8/7u0wCgYIKoZI\nzj0EAwIDSAAwRQIhAInXpLMizLP32k275s9U0d0I8x/8z4eJb/oWKe7t9OZxAiBb\nC6D+qfHM1xl/7cisyq7jZ0C5WXNn68kDVz1wiHOojg==\n-----END CERTIFICATE-----"
const PrivateKey = "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg1rtwOT3IjC/KlHW7\nTcWkRSYp3tDuD24O52mSuuh6HJyhRANCAATvx8sWfCySc8h7vNRdEbv8jSmh8iUD\nC+FBcWIHxsXbydd+GJVdq0PUfQmxkLsrftfBcRYTp9U7imIGDDHThZzc\n-----END PRIVATE KEY-----"
const TLSCertificated = "-----BEGIN CERTIFICATE-----\nMIICVzCCAf2gAwIBAgIQUeMoEOn3LSBj/DfpwHs+8jAKBggqhkjOPQQDAjB2MQsw\nCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy\nYW5jaXNjbzEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEfMB0GA1UEAxMWdGxz\nY2Eub3JnMi5leGFtcGxlLmNvbTAeFw0yMjExMDEwMjA5MDBaFw0zMjEwMjkwMjA5\nMDBaMHYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH\nEw1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcyLmV4YW1wbGUuY29tMR8wHQYD\nVQQDExZ0bHNjYS5vcmcyLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0D\nAQcDQgAEVJIP11ji+p5CHnB/KO+rj0unBTimNhGJD6YbyT72esXM/5qRIE1aTax2\n7zXc+kEOjiDh+Y5rThwgHKTAQpdHRaNtMGswDgYDVR0PAQH/BAQDAgGmMB0GA1Ud\nJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1Ud\nDgQiBCDVo0yWMnBp9i2bU1sCPiqLgFH14XXrmlqAGTJsk1kksTAKBggqhkjOPQQD\nAgNIADBFAiBYTBrMdnNMZvyOBetGxRPJCw+9+h9hDPU2IW8xP7smTAIhAOHHca19\nx/o/HCqGVBM7otFifvl18DQoh0jMihFcUJkg\n-----END CERTIFICATE-----"

func (usc *UserServiceClient) GetUserMembership(userID string) *UserMembership {
	return &UserMembership{
		MspID:          "Org2MSP",
		PeerEndpoint:   "peer0.org2.example.com:9051",
		GatewayPeer:    "peer0.org2.example.com",
		Certificate:    Certificate,
		PrivateKey:     PrivateKey,
		TLSCertificate: TLSCertificated,
	}
}

func (usc *UserServiceClient) GetUserInfo(userID string) *UserInfo {
	return &UserInfo{
		Username:       "test",
		UserID:         userID,
		OrganizationID: "Org2MSP",
		IsAdmin:        true,
	}
}
