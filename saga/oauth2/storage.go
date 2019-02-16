package oauth2

import (
	"context"
	"strings"

	"github.com/halium-project/server/resource/accesstoken"
	"github.com/halium-project/server/resource/authorizationcode"
	"github.com/halium-project/server/resource/client"
	"github.com/halium-project/server/util/errors"
	"github.com/openshift/osin"
)

const BucketName = "oauth2"

type ClientInterface interface {
	Get(ctx context.Context, cmd *client.GetCmd) (*client.Client, error)
}

type AccessTokenInterface interface {
	Create(ctx context.Context, cmd *accesstoken.CreateCmd) error
	Get(ctx context.Context, cmd *accesstoken.GetCmd) (*accesstoken.AccessToken, error)
	Delete(ctx context.Context, cmd *accesstoken.DeleteCmd) error
	FindOneByRefreshToken(ctx context.Context, cmd *accesstoken.FindOneByRefreshTokenCmd) (string, *accesstoken.AccessToken, error)
}

type AuthorizationCodeInterface interface {
	Create(ctx context.Context, cmd *authorizationcode.CreateCmd) error
	Get(ctx context.Context, cmd *authorizationcode.GetCmd) (*authorizationcode.AuthorizationCode, error)
	Delete(ctx context.Context, cmd *authorizationcode.DeleteCmd) error
}

type StorageController struct {
	client            ClientInterface
	authorizationCode AuthorizationCodeInterface
	accessToken       AccessTokenInterface
}

func NewStorageController(
	client ClientInterface,
	authorizationCode AuthorizationCodeInterface,
	accessToken AccessTokenInterface,
) *StorageController {
	return &StorageController{
		client:            client,
		authorizationCode: authorizationCode,
		accessToken:       accessToken,
	}
}

// Clone implements osin.Storage interface, but in fact does not
// clone the storage
func (t *StorageController) Clone() osin.Storage { return t }

// Close closes connections and cleans up resources
func (t *StorageController) Close() {
	// s.bucket.Close()
}

// SetClient saves client record to the storage. Client record must provide
// Id in client.Id
func (t *StorageController) SetClient(client osin.Client) error {
	return nil
}

// GetClient loads the client by id (client_id)
func (t *StorageController) GetClient(id string) (osin.Client, error) {
	client, err := t.client.Get(context.TODO(), &client.GetCmd{
		ClientID: id,
	})
	if errors.IsKind(err, errors.Validation) {
		return nil, osin.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	if client == nil {
		return nil, osin.ErrNotFound
	}

	var redirectURI string
	if len(client.RedirectURIs) > 0 {
		redirectURI = client.RedirectURIs[0]
	}

	res := osin.DefaultClient{
		Id:          id,
		Secret:      client.Secret,
		RedirectUri: redirectURI,
	}

	return &res, nil
}

// SaveAuthorize saves authorize data.
func (t *StorageController) SaveAuthorize(data *osin.AuthorizeData) error {

	err := t.authorizationCode.Create(context.TODO(), &authorizationcode.CreateCmd{
		ClientID:            data.Client.GetId(),
		Code:                data.Code,
		ExpiresIn:           int(data.ExpiresIn),
		Scopes:              strings.Split(data.Scope, ","),
		RedirectURI:         data.RedirectUri,
		State:               data.State,
		CodeChallenge:       data.CodeChallenge,
		CodeChallengeMethod: data.CodeChallengeMethod,
	})

	return err
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (t *StorageController) LoadAuthorize(code string) (*osin.AuthorizeData, error) {

	authorization, err := t.authorizationCode.Get(context.TODO(), &authorizationcode.GetCmd{
		Code: code,
	})
	if err != nil {
		return nil, err
	}

	if authorization == nil {
		return nil, osin.ErrNotFound
	}

	client, err := t.GetClient(authorization.ClientID)
	if err != nil {
		return nil, err
	}

	res := osin.AuthorizeData{
		Client:              client,
		Code:                code,
		ExpiresIn:           int32(authorization.ExpiresIn),
		Scope:               strings.Join(authorization.Scopes, ","),
		RedirectUri:         authorization.RedirectURI,
		State:               authorization.State,
		CreatedAt:           authorization.CreatedAt,
		UserData:            nil,
		CodeChallenge:       authorization.CodeChallenge,
		CodeChallengeMethod: authorization.CodeChallengeMethod,
	}

	return &res, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (t *StorageController) RemoveAuthorize(code string) error {
	err := t.authorizationCode.Delete(context.TODO(), &authorizationcode.DeleteCmd{
		Code: code,
	})

	return err
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (t *StorageController) SaveAccess(data *osin.AccessData) error {
	err := t.accessToken.Create(context.TODO(), &accesstoken.CreateCmd{
		ClientID:     data.Client.GetId(),
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    int(data.ExpiresIn),
		Scopes:       strings.Split(data.Scope, ","),
	})

	return err
}

// LoadAccess retrieves access data by token.
func (t *StorageController) LoadAccess(accessToken string) (*osin.AccessData, error) {
	token, err := t.accessToken.Get(context.TODO(), &accesstoken.GetCmd{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, osin.ErrNotFound
	}

	res := osin.AccessData{
		Client:        nil,
		AuthorizeData: nil,
		AccessData:    nil,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		ExpiresIn:     int32(token.ExpiresIn),
		Scope:         strings.Join(token.Scopes, ","),
		RedirectUri:   "",
		CreatedAt:     token.CreatedAt,
		UserData:      nil,
	}

	return &res, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (t *StorageController) RemoveAccess(accessToken string) error {
	err := t.accessToken.Delete(context.TODO(), &accesstoken.DeleteCmd{
		AccessToken: accessToken,
	})

	return err
}

// LoadRefresh retrieves refresh AccessData.
// Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (t *StorageController) LoadRefresh(refreshToken string) (*osin.AccessData, error) {
	_, session, err := t.accessToken.FindOneByRefreshToken(context.TODO(), &accesstoken.FindOneByRefreshTokenCmd{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, osin.ErrNotFound
	}

	client, err := t.GetClient(session.ClientID)
	if err != nil {
		return nil, err
	}

	res := osin.AccessData{
		Client:        client,
		AuthorizeData: nil,
		AccessData:    nil,
		AccessToken:   session.AccessToken,
		RefreshToken:  session.RefreshToken,
		ExpiresIn:     int32(session.ExpiresIn),
		Scope:         strings.Join(session.Scopes, ","),
		RedirectUri:   "",
		CreatedAt:     session.CreatedAt,
		UserData:      nil,
	}

	return &res, nil
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (t *StorageController) RemoveRefresh(refreshToken string) error {
	accessToken, _, err := t.accessToken.FindOneByRefreshToken(context.TODO(), &accesstoken.FindOneByRefreshTokenCmd{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return err
	}

	err = t.accessToken.Delete(context.TODO(), &accesstoken.DeleteCmd{
		AccessToken: accessToken,
	})

	return err
}
