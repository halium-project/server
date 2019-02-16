package authorizationcode

import "time"

const (
	Admin = "admin"
	Dev   = "dev"
)

type AuthorizationCode struct {
	// Client information
	ClientID string `json:"clientID"`

	// Token expiration in seconds
	ExpiresIn int `json:"exiresIn"`

	// Requested scope
	Scopes []string `json:"scope"`

	// Redirect Uri from request
	RedirectURI string `json:"redirectUri"`

	// State data from request
	State string `json:"state"`

	// Date created
	CreatedAt time.Time `json:"createdAt"`

	// Optional code_challenge as described in rfc7636
	CodeChallenge string `json:"codeChallenge,omitempty"`
	// Optional code_challenge_method as described in rfc7636
	CodeChallengeMethod string `json:"codeChallengeMethod,omitempty"`
}

type CreateCmd struct {
	ClientID string

	// Authorization code
	Code string

	// Token expiration in seconds
	ExpiresIn int

	// Requested scope
	Scopes []string

	// Redirect Uri from request
	RedirectURI string

	// State data from request
	State string

	// Optional code_challenge as described in rfc7636
	CodeChallenge string
	// Optional code_challenge_method as described in rfc7636
	CodeChallengeMethod string
}

type GetCmd struct {
	Code string
}

type DeleteCmd struct {
	Code string
}

var ValidAuthorizationCode = AuthorizationCode{
	ClientID:            "my-web-application",
	ExpiresIn:           10,
	Scopes:              []string{"foobar"},
	RedirectURI:         "http://some-url",
	State:               "some-ramdom-string",
	CreatedAt:           time.Now().Round(time.Millisecond),
	CodeChallenge:       "",
	CodeChallengeMethod: "",
}
