package main

type Config struct {
	SyncMinutes        int
	WebServerPort      string
	SSLCertificateFile string
	SSLPrivateKeyFile  string
	OauthClientId      string
	OauthClientSecret  string
	OauthRedirectURL   string
	ExchangeServerURL  string
	MaxResults	int
}

func (c Config) MinutesBetweenSync() int {
	return c.SyncMinutes
}

func (c Config) Port() string {
	return c.WebServerPort
}

func (c Config) Certificate() string {
	return c.SSLCertificateFile
}

func (c Config) PrivateKey() string {
	return c.SSLPrivateKeyFile
}

func (c Config) ClientId() string {
	return c.OauthClientId
}

func (c Config) ClientSecret() string {
	return c.OauthClientSecret
}

func (c Config) RedirectURL() string {
	return c.OauthRedirectURL
}

func (c Config) ExchangeURL() string {
	return c.ExchangeServerURL
}

func (c Config) MaxFetchSize() int {
	return c.MaxResults
}
