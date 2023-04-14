package redix

type XConfig struct {
	Namespace string `mapstructure:"namespace" json:"namespace"`
}

type Config struct {
	XConfig
	DSN  string `mapstructure:"dsn" json:"dsn"`
	Cert string `mapstructure:"cert" json:"cert"`
}
