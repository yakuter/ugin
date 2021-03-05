package config

type ServerConfiguration struct {
	Port                       string
	Secret                     string
	AccessTokenExpireDuration  int
	RefreshTokenExpireDuration int
}
