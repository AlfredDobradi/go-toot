package config

var (
	instanceURL         string = "https://mastodon.social"
	clientID            string = ""
	clientSecret        string = ""
	scope               string = "write:statuses"
	redirectURI         string = "urn:ietf:wg:oauth:2.0:oob"
	skipSSLVerification bool   = true
)

func InstanceURL() string {
	return instanceURL
}

func ClientID() string {
	return clientID
}

func ClientSecret() string {
	return clientSecret
}

func Scope() string {
	return scope
}

func RedirectURI() string {
	return redirectURI
}

func SkipSSLVerification() bool {
	return skipSSLVerification
}

func SetInstanceURL(url string) {
	instanceURL = url
}

func SetClientID(id string) {
	clientID = id
}

func SetClientSecret(secret string) {
	clientSecret = secret
}

func SetScope(newScope string) {
	scope = newScope
}

func SetRedirectURI(uri string) {
	redirectURI = uri
}

func SetSkipSSLVerification(val bool) {
	skipSSLVerification = val
}
