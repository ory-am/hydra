package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"

	"github.com/ory/fosite"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/x/errorsx"
)

// AccessRequestHook is called when an access token is being refreshed.
type AccessRequestHook func(ctx context.Context, requester fosite.AccessRequester) error

// TokenRefreshHookRequest is the request body sent to the token refresh hook.
//
// swagger:model tokenRefreshHookRequest
type TokenRefreshHookRequest struct {
	// Subject is the user ID of the end-user that authenticated.
	Subject string `json:"subject"`
	// GrantedScopes is the list of scopes that the end-user granted to the OAuth 2.0 client.
	GrantedScopes []string `json:"granted_scopes"`
	// GrantedAudience is the list of audiences that the end-user granted to the OAuth 2.0 client.
	GrantedAudience []string `json:"granted_audience"`
}

// TokenRefreshHookResponse is the response body sent from the token refresh hook.
//
// swagger:model tokenRefreshHookResponse
type TokenRefreshHookResponse struct {
	// Session is the session data that was returned by the hook.
	Session consent.ConsentRequestSessionData `json:"session"`
}

// TokenRefreshHook is an AccessRequestHook called for `refresh_token` grant type.
func TokenRefreshHook(config *config.Provider) AccessRequestHook {
	client := cleanhttp.DefaultPooledClient()

	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookURL := config.TokenRefreshHookURL()
		if hookURL == nil {
			return nil
		}

		if !requester.GetGrantTypes().ExactOne("refresh_token") {
			return nil
		}

		session, ok := requester.GetSession().(*Session)
		if !ok {
			return nil
		}

		reqBody := TokenRefreshHookRequest{
			Subject:         session.GetSubject(),
			GrantedScopes:   requester.GetGrantedScopes(),
			GrantedAudience: requester.GetGrantedAudience(),
		}
		reqBodyBytes, err := json.Marshal(&reqBody)
		if err != nil {
			return errorsx.WithStack(err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, hookURL.String(), bytes.NewReader(reqBodyBytes))
		if err != nil {
			return errorsx.WithStack(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return errorsx.WithStack(err)
		}
		if resp.StatusCode != http.StatusOK {
			return errorsx.WithStack(fosite.ErrServerError.WithDebugf("Token refresh hook responded with %s status.", resp.Status))
		}
		defer resp.Body.Close()

		var respBody TokenRefreshHookResponse
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return errorsx.WithStack(err)
		}

		session.Extra = respBody.Session.AccessToken
		idTokenClaims := session.IDTokenClaims()
		idTokenClaims.Extra = respBody.Session.IDToken

		return nil
	}
}
