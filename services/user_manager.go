package services

type UserMembership struct {
	MspID          string `mapstructure:"msp_id"`
	PeerEndpoint   string `mapstructure:"peer_endpoint"`
	GatewayPeer    string `mapstructure:"gateway_peer"`
	Certificate    string `mapstructure:"certificate_path"`
	PrivateKey     string `mapstructure:"private_key_path"`
	TLSCertificate string `mapstructure:"tls_certificate_path"`
}

type UserManager struct {
}

func NewUserManager() *UserManager {
	return new(UserManager)
}

const Certificate = "-----BEGIN CERTIFICATE-----\nMIICKjCCAdCgAwIBAgIRAIsc83pQ/7unOwJAufXnMO0wCgYIKoZIzj0EAwIwczEL\nMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG\ncmFuY2lzY28xGTAXBgNVBAoTEG9yZzIuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh\nLm9yZzIuZXhhbXBsZS5jb20wHhcNMjIxMDI5MDQ0OTAwWhcNMzIxMDI2MDQ0OTAw\nWjBrMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN\nU2FuIEZyYW5jaXNjbzEOMAwGA1UECxMFYWRtaW4xHzAdBgNVBAMMFkFkbWluQG9y\nZzIuZXhhbXBsZS5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASjDFIiZd1U\nh+cNfwpXoTuyS0JhOc76mRsKhv7braglmjIyP8KYy4EpijwrSaBx9xaQPPClWhxL\npLcn06KUV49Eo00wSzAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADArBgNV\nHSMEJDAigCAmI382PId7ktnuYDNKCy3eZcbADJfAfHv3XiMunF8kBzAKBggqhkjO\nPQQDAgNIADBFAiEA0OUw+p6mX4o6jYtJEo+nwjO8Mmd0jQO9acLsAPoe/KMCICYg\nwG8S+9HrDvnUx1ZuIb6oPXATn5oEewQcjoGsw835\n-----END CERTIFICATE-----"
const PrivateKey = "-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgMfo7dYiTofm2l4B1\nziMkxFM8wOT3wF1UB1cE4iWXm2ihRANCAASjDFIiZd1Uh+cNfwpXoTuyS0JhOc76\nmRsKhv7braglmjIyP8KYy4EpijwrSaBx9xaQPPClWhxLpLcn06KUV49E\n-----END PRIVATE KEY-----"
const TLSCertificated = "-----BEGIN CERTIFICATE-----\nMIICVjCCAf2gAwIBAgIQWaL39vKO341L0ZausEwtATAKBggqhkjOPQQDAjB2MQsw\nCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy\nYW5jaXNjbzEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEfMB0GA1UEAxMWdGxz\nY2Eub3JnMi5leGFtcGxlLmNvbTAeFw0yMjEwMjkwNDQ5MDBaFw0zMjEwMjYwNDQ5\nMDBaMHYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH\nEw1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcyLmV4YW1wbGUuY29tMR8wHQYD\nVQQDExZ0bHNjYS5vcmcyLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0D\nAQcDQgAEdEK0sl/IIIf9ZDiqnw2QbHoDRk79dxqff9igm0MUmE7vXK6Yx3ptvXk8\n4JrRHoJj6UD0y4NyLWiw/p0D7XTMR6NtMGswDgYDVR0PAQH/BAQDAgGmMB0GA1Ud\nJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1Ud\nDgQiBCAxof2Vwl/pVj6HjArcu7GtV/WzaM/mn1/CRlg+uq7QtDAKBggqhkjOPQQD\nAgNHADBEAiAHTtxNiiQhbyKKCi2T27OJqh0MUTYr0xWotAs8u/PfeQIgKVHg64vZ\nceBKZSmwQyU8hx7QrnBptQFK780IJAKXV3I=\n-----END CERTIFICATE-----"

func (userManager *UserManager) GetUserMembership(userID string) *UserMembership {
	return &UserMembership{
		MspID:          "Org2MSP",
		PeerEndpoint:   "localhost:9051",
		GatewayPeer:    "peer0.org2.example.com",
		Certificate:    Certificate,
		PrivateKey:     PrivateKey,
		TLSCertificate: TLSCertificated,
	}
}
