package config

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"crypto/sha256"

	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/fosite/token/hmac"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	r "gopkg.in/dancannon/gorethink.v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// These are used by client commands
	ClusterURL          string `mapstructure:"CLUSTER_URL" yaml:"cluster_url"`
	ClientID            string `mapstructure:"CLIENT_ID" yaml:"client_id,omitempty"`
	ClientSecret        string `mapstructure:"CLIENT_SECRET" yaml:"client_secret,omitempty"`

	// These are used by the host command
	BindPort            int `mapstructure:"PORT" yaml:"-"`
	BindHost            string `mapstructure:"HOST" yaml:"-"`
	Issuer              string `mapstructure:"ISSUER" yaml:"-"`
	SystemSecret        string `mapstructure:"SYSTEM_SECRET" yaml:"-"`
	DatabaseURL         string `mapstructure:"DATABASE_URL" yaml:"-"`
	ConsentURL          string `mapstructure:"CONSENT_URL" yaml:"-"`
	BCryptWorkFactor    int `mapstructure:"BCRYPT_COST" yaml:"-"`
	AccessTokenLifespan string `mapstructure:"ACCESS_TOKEN_LIFESPAN" yaml:"-"`
	AuthCodeLifespan    string `mapstructure:"AUTH_CODE_LIFESPAN" yaml:"-"`
	IDTokenLifespan     string `mapstructure:"ID_TOKEN_LIFESPAN" yaml:"-"`
	ChallengeTokenLifespan   string `mapstructure:"CHALLENGE_TOKEN_LIFESPAN" yaml:"-"`
	ForceHTTP           bool `yaml:"-"`

	cluster             *url.URL `yaml:"-"`
	oauth2Client        *http.Client `yaml:"-"`
	context             *Context `yaml:"-"`
	sync.Mutex `yaml:"-"`
}

func (c *Config) GetChallengeTokenLifespan() time.Duration {
	d, err := time.ParseDuration(c.ChallengeTokenLifespan)
	if err != nil {
		logrus.Warnf("Could not parse challenge token lifespan value (%s). Defaulting to 10m", c.AccessTokenLifespan)
		return time.Minute * 10
	}
	return d
}

func (c *Config) GetAccessTokenLifespan() time.Duration {
	d, err := time.ParseDuration(c.AccessTokenLifespan)
	if err != nil {
		logrus.Warnf("Could not parse access token lifespan value (%s). Defaulting to 1h", c.AccessTokenLifespan)
		return time.Hour
	}
	return d
}

func (c *Config) GetAuthCodeLifespan() time.Duration {
	d, err := time.ParseDuration(c.AuthCodeLifespan)
	if err != nil {
		logrus.Warnf("Could not parse auth code lifespan value (%s). Defaulting to 10m", c.AuthCodeLifespan)
		return time.Minute * 10
	}
	return d
}

func (c *Config) GetIDTokenLifespan() time.Duration {
	d, err := time.ParseDuration(c.IDTokenLifespan)
	if err != nil {
		logrus.Warnf("Could not parse id token lifespan value (%s). Defaulting to 1h", c.IDTokenLifespan)
		return time.Hour
	}
	return d
}

func (c *Config) Context() *Context {
	if c.context != nil {
		return c.context
	}

	var connection interface{} = &MemoryConnection{}
	if c.DatabaseURL != "" {
		u, err := url.Parse(c.DatabaseURL)
		if err != nil {
			logrus.Fatalf("Could not parse DATABASE_URL: %s", err)
		}

		switch u.Scheme {
		case "rethinkdb":
			connection = &RethinkDBConnection{URL: u}
			break
		default:
			logrus.Fatalf("Unkown DSN in DATABASE_URL: %s", c.DatabaseURL)
		}
	}

	var manager ladon.Manager
	switch con := connection.(type) {
	case *MemoryConnection:
		logrus.Printf("DATABASE_URL not set, connecting to ephermal in-memory database.")
		manager = ladon.NewMemoryManager()
		break
	case *RethinkDBConnection:
		logrus.Printf("DATABASE_URL set, connecting to RethinkDB.")
		con.CreateTableIfNotExists("hydra_policies")
		m := &ladon.RethinkManager{
			Session: con.GetSession(),
			Table:   r.Table("hydra_policies"),
		}
		m.Watch(context.Background())
		if err := m.ColdStart(); err != nil {
			logrus.Fatalf("Could not fetch initial state: %s", err)
		}
		manager = m
		break
	default:
		panic("Unknown connection type.")
	}

	c.context = &Context{
		Connection: connection,
		Hasher: &hash.BCrypt{
			WorkFactor: c.BCryptWorkFactor,
		},
		LadonManager: manager,
		FositeStrategy: &strategy.HMACSHAStrategy{
			Enigma: &hmac.HMACStrategy{
				GlobalSecret: c.GetSystemSecret(),
			},
			AccessTokenLifespan:   c.GetAccessTokenLifespan(),
			AuthorizeCodeLifespan: c.GetAuthCodeLifespan(),
		},
	}

	return c.context
}

func (c *Config) Resolve(join ...string) *url.URL {
	if c.cluster == nil {
		cluster, err := url.Parse(c.ClusterURL)
		c.cluster = cluster
		pkg.Must(err, "Could not parse cluster url: %s", err)
	}

	if len(join) == 0 {
		return c.cluster
	}

	return pkg.JoinURL(c.cluster, join...)
}

func (c *Config) OAuth2Client(cmd *cobra.Command) *http.Client {
	if c.oauth2Client != nil {
		return c.oauth2Client
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     pkg.JoinURLStrings(c.ClusterURL, "/oauth2/token"),
		Scopes: []string{
			"core",
			"hydra",
		},
	}

	ctx := context.Background()
	if ok, _ := cmd.Flags().GetBool("skip-tls-verify"); ok {
		fmt.Println("Warning: Skipping TLS Certificate Verification.")
		ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}})
	}

	_, err := oauthConfig.Token(ctx)
	if err != nil {
		fmt.Printf("Could not authenticate, because: %s\n", err)
		fmt.Println("Did you forget to log on? Run `hydra connect`.")
		fmt.Println("Did you run Hydra without a valid TLS certificate? Make sure to use the `--skip-tls-verify` flag.")
		fmt.Println("Did you know you can skip `hydra connect` when running `hydra host --dangerous-auto-logon`? DO NOT use this flag in production!")
		os.Exit(1)
	}

	c.oauth2Client = oauthConfig.Client(ctx)
	return c.oauth2Client
}

func (c *Config) GetSystemSecret() []byte {
	var secret = []byte(c.SystemSecret)
	if len(secret) >= 16 {
		hash := sha256.Sum256(secret)
		secret = hash[:]
		c.SystemSecret = string(secret)
		return secret
	}

	logrus.Warnf("Expected system secret to be at least %d characters long, got %d characters.", 32, len(c.SystemSecret))
	logrus.Infoln("Generating a random system secret...")
	var err error
	secret, err = pkg.GenerateSecret(32)
	pkg.Must(err, "Could not generate global secret: %s", err)
	logrus.Infof("Generated system secret: %s", secret)
	hash := sha256.Sum256(secret)
	secret = hash[:]
	c.SystemSecret = string(secret)
	logrus.Warnln("WARNING: DO NOT generate system secrets in production. The secret will be leaked to the logs.")
	return secret
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}

func (c *Config) Persist() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return errors.New(err)
	}

	if err := ioutil.WriteFile(viper.ConfigFileUsed(), out, 0700); err != nil {
		return errors.Errorf(`Could not write to "%s" because: %s`, viper.ConfigFileUsed(), err)
	}

	return nil
}
