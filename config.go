package redix

type Config struct {
	XConfig `mapstructure:",squash"`
	DSN     string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
	Cert    string `json:"cert" yaml:"cert" mapstructure:"cert"`
}

type XConfig struct {
	Namespace string `json:"namespace" yaml:"namespace" mapstructure:"namespace"`
}

func (c *XConfig) AppendNamespace(ns string) {
	if c.Namespace == "" {
		c.Namespace = ns
		return
	}

	c.Namespace = c.Namespace + keyDelimiter + ns
}
