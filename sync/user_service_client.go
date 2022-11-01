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

const Certificate = "-----BEGIN CERTIFICATE-----\nMIICKjCCAdGgAwIBAgIRAKMjcDdKohrWWy/OK0CDbTUwCgYIKoZIzj0EAwIwczEL\nMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG\ncmFuY2lzY28xGTAXBgNVBAoTEG9yZzIuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh\nLm9yZzIuZXhhbXBsZS5jb20wHhcNMjIxMDI5MjAxNzAwWhcNMzIxMDI2MjAxNzAw\nWjBsMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN\nU2FuIEZyYW5jaXNjbzEPMA0GA1UECxMGY2xpZW50MR8wHQYDVQQDDBZVc2VyMUBv\ncmcyLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAErmM4OXNl\nMQeF/7gBBiCtsuYlpTmUjIKcAOJr2N26U2was0E1zBywghpgDdfS/sXP1NAGNF8U\n8h3IkpF4jJ0l2aNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYD\nVR0jBCQwIoAgZZVyXUA7KgDZxjd7m7Oa1QphNLDUYqMdQxwRVCD/0vwwCgYIKoZI\nzj0EAwIDRwAwRAIgUYe8sgZtmceDygkC5sK6GF7OsXTF6w6rrLjHVEBYaXACIGKa\n8z/a7maOEQD7Nc2vjAtABczmigd2RM/MkkUfSuSH\n-----END CERTIFICATE-----"
const PrivateKey = "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgLzIeq9aqObE2UW9G\nemfbexvMUjVkVJz+mSb9Xob1TXWhRANCAASuYzg5c2UxB4X/uAEGIK2y5iWlOZSM\ngpwA4mvY3bpTbBqzQTXMHLCCGmAN19L+xc/U0AY0XxTyHciSkXiMnSXZ\n-----END PRIVATE KEY-----"
const TLSCertificated = "-----BEGIN CERTIFICATE-----\nMIICWTCCAf6gAwIBAgIRAM+WdVUdRymcDcxpjVh7YtUwCgYIKoZIzj0EAwIwdjEL\nMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG\ncmFuY2lzY28xGTAXBgNVBAoTEG9yZzIuZXhhbXBsZS5jb20xHzAdBgNVBAMTFnRs\nc2NhLm9yZzIuZXhhbXBsZS5jb20wHhcNMjIxMDI5MjAxNzAwWhcNMzIxMDI2MjAx\nNzAwWjB2MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UE\nBxMNU2FuIEZyYW5jaXNjbzEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEfMB0G\nA1UEAxMWdGxzY2Eub3JnMi5leGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49\nAwEHA0IABIMEAEp/7fuaoRYfT+NJxtRo+eEhfEIdtmFgo/a64TjDpjzpDRUR6gl+\nKl0d0ORdfc9YuBhMDfvthSVRwg3/gXSjbTBrMA4GA1UdDwEB/wQEAwIBpjAdBgNV\nHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zApBgNV\nHQ4EIgQgotLmyBZEOB5f4asV4aq6oT9jOox8YvsbqGsPrhlGCgAwCgYIKoZIzj0E\nAwIDSQAwRgIhAO2Unk65qBg+wvnBri1/wzkjNT0IdqDKf1L9Jw0O8rYhAiEAjMjW\n6kf69yPtgfFJ42tM5ds30BwCotVxSekYQQSAut8=\n-----END CERTIFICATE-----"

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
