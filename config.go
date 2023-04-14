package redix

type Config struct {
	XConfig `mapstructure:",squash"`
	DSN     string `mapstructure:"dsn" json:"dsn"`
	Cert    string `mapstructure:"cert" json:"cert"`
}

type XConfig struct {
	Namespace string `mapstructure:"namespace" json:"namespace"`
}

func (c *XConfig) AppendNamespace(ns string) {
	if c.Namespace == "" {
		c.Namespace = ns
		return
	}

	c.Namespace = c.Namespace + keyDelimiter + ns
}
