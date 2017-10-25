package jwk_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
)

var managers = map[string]Manager{}

var testGenerators = (&Handler{}).GetGenerators()

var ts *httptest.Server
var httpManager *HTTPManager

func init() {
	localWarden, httpClient := compose.NewMockFirewall(
		"tests",
		"alice",
		fosite.Arguments{
			"hydra.keys.create",
			"hydra.keys.get",
			"hydra.keys.delete",
			"hydra.keys.update",
		}, &ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:keys:<RS256|ES256|ES521|HS256><faz|bar|foo|anonymous><.*>"},
			Actions:   []string{"create", "get", "delete", "update"},
			Effect:    ladon.AllowAccess,
		}, &ladon.DefaultPolicy{
			ID:        "2",
			Subjects:  []string{"alice", ""},
			Resources: []string{"rn:hydra:keys:<RS256|ES256|ES521|HS256>anonymous<.*>"},
			Actions:   []string{"get"},
			Effect:    ladon.AllowAccess,
		},
	)

	router := httprouter.New()
	h := Handler{
		Manager: &MemoryManager{},
		W:       localWarden,
		H:       herodot.NewJSONWriter(nil),
	}
	h.SetRoutes(router)
	ts := httptest.NewServer(router)
	u, _ := url.Parse(ts.URL + "/keys")
	managers["memory"] = &MemoryManager{}
	httpManager = &HTTPManager{Client: httpClient, Endpoint: u}
	managers["http"] = httpManager
}

var encryptionKey, _ = RandomBytes(32)

func TestMain(m *testing.M) {
	connectToPG()
	connectToMySQL()

	s := m.Run()
	integration.KillAll()
	os.Exit(s)
}

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	managers["postgres"] = s
}

func connectToMySQL() {
	var db = integration.ConnectToMySQL()
	s := &SQLManager{DB: db, Cipher: &AEAD{Key: encryptionKey}}

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	managers["mysql"] = s
}

func TestHTTPManagerPublicKeyGet(t *testing.T) {
	for algo, testGenerator := range testGenerators {
		if algo == "HS256" {
			// this is a symmetrical algorithm
			continue
		}

		anonymous := &HTTPManager{Endpoint: httpManager.Endpoint, Client: http.DefaultClient}
		ks, _ := testGenerator.Generate("")
		priv := ks.Key("private")

		name := "http"
		m := httpManager

		_, err := m.GetKey(algo+"anonymous", "baz")
		pkg.AssertError(t, true, err, name)

		err = m.AddKey(algo+"anonymous", First(priv))
		pkg.AssertError(t, false, err, name)

		time.Sleep(time.Millisecond * 100)

		got, err := anonymous.GetKey(algo+"anonymous", "private")
		pkg.RequireError(t, false, err, name)
		assert.Equal(t, priv, got.Keys, "%s", name)
	}
}

func TestManagerKey(t *testing.T) {
	for algo, testGenerator := range testGenerators {
		if algo == "HS256" {
			// this is a symmetrical algorithm
			continue
		}

		ks, err := testGenerator.Generate("")
		if err != nil {
			t.Fatal(err)
		}

		for name, m := range managers {
			t.Run(fmt.Sprintf("case=%s/%s", algo, name), func(t *testing.T) {
				TestHelperManagerKey(m, algo, ks)(t)
			})
		}

		priv := ks.Key("private")
		err = managers["http"].AddKey("nonono", First(priv))
		assert.NotNil(t, err)
	}
}

func TestManagerKeySet(t *testing.T) {
	for algo, testGenerator := range testGenerators {
		ks, err := testGenerator.Generate("")
		if err != nil {
			t.Fatal(err)
		}
		ks.Key("private")

		for name, m := range managers {
			t.Run(fmt.Sprintf("case=%s/%s", algo, name), func(t *testing.T) {
				TestHelperManagerKeySet(m, algo, ks)(t)
			})
		}

		err = managers["http"].AddKeySet("nonono", ks)
		assert.NotNil(t, err)
	}
}
