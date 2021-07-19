package cli

import (
	"context"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_toImportJSONWebKey(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)
	router := x.NewRouterPublic()

	h := reg.KeyHandler()
	m := reg.KeyManager()

	h.SetRoutes(router.RouterAdmin(), router, func(h http.Handler) http.Handler {
		return h
	})
	testServer := httptest.NewServer(router)

	cmd := cobra.Command{
		Use: "key",
	}
	cmd.Flags().String("use", "sig", "Sets the \"use\" value of the JSON Web Key if not \"use\" value was defined by the key itself")
	cmd.Flags().Bool("fake-tls-termination", false, "Sets the \"use\" value of the JSON Web Key if not \"use\" value was defined by the key itself")
	cmd.Flags().String("access-token", "", "Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN")
	cmd.Flags().String("endpoint", "", "Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL. A unix socket can be set in the form unix:///path/to/socket")
	cmd.Flags().Bool("skip-tls-verify", true, "Foolishly accept TLS certificates signed by unknown certificate authorities")
	cmd.Flags().String("default-key-id","cae6b214-fb1e-4ebc-9019-95286a62eabc","A fallback value for keys without \"kid\" attribute to be stored with a common \"kid\", e.g. private/public keypairs")
	os.Setenv("HYDRA_URL", testServer.URL)

	NewHandler().Keys.ImportKeys(&cmd, []string{"import-1", "../test/private_key.json", "../test/public_key.json"})
	//running again to make sure the row in storage is not deleted issue: #2436
	NewHandler().Keys.ImportKeys(&cmd, []string{"import-1", "../test/private_key.json", "../test/public_key.json"})
	v, _ := m.GetKeySet(context.TODO(), "import-1")
	assert.Equal(t, len(v.Keys),2)

	NewHandler().Keys.ImportKeys(&cmd, []string{"import-1", "../test/private_key.json", "../test/another_public_key.json"})
	v, _ = m.GetKeySet(context.TODO(), "import-1")
	assert.Equal(t, len(v.Keys),3)

	NewHandler().Keys.ImportKeys(&cmd, []string{"import-2", "../test/private_key.json", "../test/public_key.json"})
	v, _ = m.GetKeySet(context.TODO(), "import-2")
	assert.NotEmpty(t, v)

}

func Test_toUpdateKeyList(t *testing.T) {

	dummyPEM := []byte(`
		-----BEGIN PUBLIC KEY-----
		1MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAPf64dykufSkwnvUiBAwd5Si0K6t4m5i
		qJD8TmLJCmFjKaOUa6nszcFt/FkAuORfdlrD9mEZLPrPx74RSluyTBMCAwEAAQ==
		-----END PUBLIC KEY-----
	`)

	anotherPEM := []byte(`
		-----BEGIN PUBLIC KEY-----
		2MFwwDAQYJKoZIhvcNAQEBBQADSwAwSAJBAPf64dykufSkwnvUiBAwd5Si0K6t4m5i
		qJD8TmLJCmFjKaOUa6nszcFt/FkAuORfdlrD9mEZLPrPx74RSluyTBMCAwEAAQ==
		-----END PUBLIC KEY-----
	`)

	var set jose.JSONWebKeySet
	set.Keys = updateKey(set,toSDKFriendlyJSONWebKey(dummyPEM, "static-sid", "sig") )
	set.Keys = updateKey(set,toSDKFriendlyJSONWebKey(dummyPEM, "another-sid", "sig") )
	assert.Equal(t, len(set.Keys),2)

	//overriding static-sid key with publicAnotherPEM
	set.Keys = updateKey(set,toSDKFriendlyJSONWebKey(anotherPEM, "static-sid", "sig") )
	assert.Equal(t, len(set.Keys),2)

	key := set.Key("static-sid")[0]
	verifyKey := toSDKFriendlyJSONWebKey(anotherPEM, "static-sid", "sig")
	assert.Equal(t, key,verifyKey)

	//adding a third key with different sid
	set.Keys = updateKey(set,toSDKFriendlyJSONWebKey(anotherPEM, "different-sid", "sig") )
	assert.Equal(t, len(set.Keys),3)



}
