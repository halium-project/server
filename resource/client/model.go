package client

// Client provides the underlying structured make up of an OAuth2.0 Client.
//
// In order to update mongo records efficiently
// omitempty is used for all bson casting, with exception to ID, as this should always be provided in queries and
// updates.
type Client struct {
	// Client unique identifier.
	//
	// It is human readable and defined by the user, which should be the store
	// in most of the cases.
	ID string `json:"id"`

	// Name is the human-readable string name of the client to be presented to the
	// end-user during authorization.
	Name string `json:"name"`

	// Secret is the client's secret. The secret will be included in the create request as cleartext, and then
	// never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users
	// that they need to write the secret down as it will not be made available again.
	Secret string `json:"secret,omitempty"`

	// RedirectURIs is an array of allowed redirect urls for the client, for example:
	// http://mydomain/oauth/callback.
	RedirectURIs []string `json:"redirectURIs"`

	// GrantTypes is an array of grant types the client is allowed to use.
	//
	// Pattern: client_credentials|authorize_code|implicit|refresh_token|password
	GrantTypes []string `json:"grantTypes"`

	// ResponseTypes is an array of the OAuth 2.0 response type strings that the client can
	// use at the authorization endpoint.
	//
	// Pattern: id_token|code|token
	ResponseTypes []string `json:"responseTypes"`

	// Scope is a string containing a space-separated list of scope values (as
	// described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
	// can use when requesting access tokens.
	//
	// Pattern: ([a-zA-Z0-9\.]+\s)+
	Scopes []string `json:"scopes"`

	// Public is a boolean that identifies this client as public, meaning that it
	// does not have a secret. It will disable the client_credentials grant type for this client if set.
	Public bool `json:"public"`
}

type GetAllCmd struct{}

type GetCmd struct {
	ClientID string
}

type DeleteCmd struct {
	ClientID string
}

type ValidateCmd struct {
	ClientID     string
	ClientSecret string
}

type GetByNameCmd struct {
	Name string
}

type CreateCmd struct {
	ID            string
	Name          string
	RedirectURIs  []string
	GrantTypes    []string
	ResponseTypes []string
	Scopes        []string
	Public        bool
}

var ValidClient = Client{
	ID:            "my-web-application",
	Name:          "My Web Application",
	Secret:        "some-hashed-secret",
	RedirectURIs:  []string{"http://mydomain/oauth/callback"},
	GrantTypes:    []string{"client_credentials", "authorize_code", "implicit", "refresh_token", "password"},
	ResponseTypes: []string{"code", "token"},
	Scopes:        []string{"user", "admin"},
	Public:        false,
}
