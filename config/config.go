package config

type Config struct {
	Username      string
	SecurityCode1 string
	SecurityCode2 string
	Source        string
	URL           string
}

func New(username, securityCode1, securityCode2 string) *Config {
	return &Config{
		Username:      username,
		SecurityCode1: securityCode1,
		SecurityCode2: securityCode2,
		Source:        "sallandpioneers/go-eboekhouden api",
		URL:           "https://soap.e-boekhouden.nl/soap.asmx?WSDL",
	}
}

func (cfg *Config) WithSource(source string) *Config {
	cfg.Source = source
	return cfg
}

func (cfg *Config) WithURL(url string) *Config {
	cfg.URL = url
	return cfg
}
