package autoload

type JwtConfig struct {
	Secret       string `mapstructure:"secret"`
	HeaderPrefix string `mapstructure:"header_prefix"`
	Expiration   int    `mapstructure:"expiration"`
	RefreshTTL   int    `mapstructure:"refresh_time"`
}
