package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/connection"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/policy"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
)

type Handler struct {
	Clients     *client.Handler
	Connections *connection.Handler
	Keys        *jwk.Handler
	OAuth2      *oauth2.Handler
	Policy      *policy.Handler
	Warden      *warden.WardenHandler
}

func (h *Handler) Start(c *config.Config, router *httprouter.Router) {
	ctx := c.Context()

	// Set up warden
	clientsManager := newClientManager(c)
	injectFositeStore(c, clientsManager)
	ctx.Warden = &warden.LocalWarden{
		Warden: &ladon.Ladon{
			Manager: ctx.LadonManager,
		},
		TokenValidator: &core.CoreValidator{
			AccessTokenStrategy: ctx.FositeStrategy,
			AccessTokenStorage:  ctx.FositeStore,
		},
		Issuer: c.Issuer,
	}

	// Set up handlers
	h.Clients = newClientHandler(c, router, clientsManager)
	h.Keys = newJWKHandler(c, router)
	h.Connections = newConnectionHandler(c, router)
	h.Policy = newPolicyHandler(c, router)
	h.OAuth2 = newOAuth2Handler(c, router, h.Keys.Manager)
	h.Warden = warden.NewHandler(c, router)

	// Create root account if new install
	h.createRS256KeysIfNotExist(c, oauth2.ConsentEndpointKey, "private")
	h.createRS256KeysIfNotExist(c, oauth2.ConsentChallengeKey, "private")

	h.createRootIfNewInstall(c)
}

func (h *Handler) createRS256KeysIfNotExist(c *config.Config, set, lookup string) {
	ctx := c.Context()
	generator := jwk.RS256Generator{}

	if _, err := ctx.KeyManager.GetKey(set, lookup); errors.Is(err, pkg.ErrNotFound) {
		logrus.Warnf("Key pair for signing %s is missing. Creating new one.", set)

		keys, err := generator.Generate("")
		pkg.Must(err, "Could not generate %s key: %s", set, err)

		err = ctx.KeyManager.AddKeySet(set, keys)
		pkg.Must(err, "Could not persist %s key: %s", set, err)
	}
}

func (h *Handler) createRootIfNewInstall(c *config.Config) {
	ctx := c.Context()

	clients, err := h.Clients.Manager.GetClients()
	pkg.Must(err, "Could not fetch client list: %s", err)
	if len(clients) != 0 {
		return
	}

	rs, err := pkg.GenerateSecret(16)
	pkg.Must(err, "Could notgenerate secret because %s", err)
	secret := string(rs)

	logrus.Warn("No clients were found. Creating a temporary root client...")
	root := &client.Client{
		Name:          "This temporary client is generated by hydra and is granted all of hydra's administrative privileges. It must be removed when everything is set up.",
		GrantTypes:    []string{"client_credentials", "authorization_code"},
		ResponseTypes: []string{"token", "code"},
		GrantedScopes: []string{"hydra", "core"},
		RedirectURIs:  []string{"http://localhost:4445/callback"},
		Secret:        secret,
	}

	err = h.Clients.Manager.CreateClient(root)
	pkg.Must(err, "Could not create temporary root because %s", err)
	err = ctx.LadonManager.Create(&ladon.DefaultPolicy{
		Description: "This is a policy created by hydra and issued to the first client. It grants all of hydra's administrative privileges to the client and enables the client_credentials response type.",
		Subjects:    []string{root.GetID()},
		Effect:      ladon.AllowAccess,
		Resources:   []string{"rn:hydra:<.*>"},
		Actions:     []string{"<.*>"},
	})
	pkg.Must(err, "Could not create admin policy because %s", err)

	c.Lock()
	c.ClientID = root.ID
	c.ClientSecret = string(secret)
	c.Unlock()

	logrus.Warn("Temporary root client created.")
	logrus.Warnf("client_id: %s", root.GetID())
	logrus.Warnf("client_secret: %s", string(secret))
	logrus.Warn("The root client must be removed in production. The root's credentials could be accidentally logged.")
}
