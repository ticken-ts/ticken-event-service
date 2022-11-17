package config

type DevConfig struct {
	JWTPublicKey  string `mapstructure:"jwt_public_key"`
	JWTPrivateKey string `mapstructure:"jwt_private_key"`
}
