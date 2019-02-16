package accesstoken

import "time"

const (
	Admin = "admin"
	Dev   = "dev"
)

type AccessToken struct {
	// Client information
	ClientID string `json:"clientID"`

	// AccessToken token
	AccessToken string `json:"accessToken"`

	// Refresh Token. Can be blank
	RefreshToken string `json:"refreshToken"`

	// Token expiration in seconds
	ExpiresIn int `json:"expiresIn"`

	// Requested scope
	Scopes []string `json:"scopes"`

	// Date created
	CreatedAt time.Time `json:"createdAt"`
}

type CreateCmd struct {
	ClientID     string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Scopes       []string
}

type GetCmd struct {
	AccessToken string
}

type FindOneByRefreshTokenCmd struct {
	RefreshToken string
}

type DeleteCmd struct {
	AccessToken string
}

var ValidAccessToken = AccessToken{
	ClientID:     "658fae6d-f9ae-4c5c-8121-bac88ff2ee4b",
	AccessToken:  "some-access-token",
	RefreshToken: "some-refresh-token",
	ExpiresIn:    3600,
	Scopes:       []string{"users", "foobar"},
	CreatedAt:    time.Now().Round(time.Millisecond),
}
