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

const Certificate = "-----BEGIN CERTIFICATE-----\nMIICKzCCAdGgAwIBAgIRAK5bhgUp7du2OGga8+xJJtowCgYIKoZIzj0EAwIwczEL\nMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG\ncmFuY2lzY28xGTAXBgNVBAoTEG9yZzIuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh\nLm9yZzIuZXhhbXBsZS5jb20wHhcNMjIxMTA2MDM0MjAwWhcNMzIxMTAzMDM0MjAw\nWjBsMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN\nU2FuIEZyYW5jaXNjbzEPMA0GA1UECxMGY2xpZW50MR8wHQYDVQQDDBZVc2VyMUBv\ncmcyLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzByBoX+G\ncsbEOgNU0haHTitq3RVYRVnTEWucHLh0S2UB4Z4x7rQ1zlxWyVvZAWBf2PRF9jtZ\nVngdS21T1j5rqqNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYD\nVR0jBCQwIoAgjqFSt26biGZ9coQbcbMtL5XaGEI2+ZhqBzBzVUrmieYwCgYIKoZI\nzj0EAwIDSAAwRQIhAPxvPtQtCFJBbv1wOagQ9FHcXQTzrPp7mCcoV01lAxyJAiBs\n0P6SJcSaie+WSLXH4wHXHTi5SLbuQi09FVLHOqqTZw==\n-----END CERTIFICATE-----"
const PrivateKey = "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg4PkHwQmtGPYvLZny\nKJ315+4sADziI+FOKAq+hDffcNmhRANCAATMHIGhf4ZyxsQ6A1TSFodOK2rdFVhF\nWdMRa5wcuHRLZQHhnjHutDXOXFbJW9kBYF/Y9EX2O1lWeB1LbVPWPmuq\n-----END PRIVATE KEY-----"
const TLSCertificated = "-----BEGIN CERTIFICATE-----\nMIICVzCCAf2gAwIBAgIQDQ6U4LXdsXQA8nIg23yp4zAKBggqhkjOPQQDAjB2MQsw\nCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy\nYW5jaXNjbzEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEfMB0GA1UEAxMWdGxz\nY2Eub3JnMi5leGFtcGxlLmNvbTAeFw0yMjExMDYwMzQyMDBaFw0zMjExMDMwMzQy\nMDBaMHYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH\nEw1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcyLmV4YW1wbGUuY29tMR8wHQYD\nVQQDExZ0bHNjYS5vcmcyLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0D\nAQcDQgAEWaVQecbPSM5YcZBLwFx/LYKBorMYs10yK2MQKWRyxbkxvE1RNFASEcL/\nRp/CRicOdfKyEUdr/B54SfajeaimxKNtMGswDgYDVR0PAQH/BAQDAgGmMB0GA1Ud\nJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1Ud\nDgQiBCDpVBXrFEunx729EG7OYhw4Adnbs9BdnGYaiefP/7NQnjAKBggqhkjOPQQD\nAgNIADBFAiEAvbagbf7ytLJ2MVdCUzI4Pojw6QUTqNRVcQB4epNScKoCIH4XZhuj\n2Nm/x9Sgs2HVOWZQkkJSSlR0aEoAXxXJpN4w\n-----END CERTIFICATE-----"

func (usc *UserServiceClient) GetUserMembership(userID string) *UserMembership {
	return &UserMembership{
		MspID:          "Org2MSP",
		PeerEndpoint:   "localhost:9051",
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
