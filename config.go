package main

type Config struct {
	SyncInterval          string
	WebServerPort         string
	SSLCertificateFile    string
	SSLPrivateKeyFile     string
	OauthClientId         string
	OauthClientSecret     string
	OauthRedirectURL      string
	ExchangeServerURL     string
	ExchangeServerVersion string
	UserExchangeDomain    string
	CalendarLookAheadDays int
	MaxResults            int
}

func (c Config) TimeBetweenSync() string {
	return c.SyncInterval
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

func (c Config) ExchangeVersion() string {
	return c.ExchangeServerVersion
}

func (c Config) UserDomain() string {
	return c.UserExchangeDomain
}

func (c Config) LookAheadDays() int {
	return c.CalendarLookAheadDays
}

func (c Config) MaxFetchSize() int {
	return c.MaxResults
}
