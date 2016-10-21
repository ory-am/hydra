package oauth2

import (
	"bytes"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type HTTPRecovator struct {
	Client   *http.Client
	Dry      bool
	Endpoint *url.URL
}

func (r *HTTPRecovator) SetClient(c *clientcredentials.Config) {
	r.Client = c.Client(oauth2.NoContext)
}

func (r *HTTPRecovator) RevokeToken(ctx context.Context, token string) (error) {
	var ep = *r.Endpoint
	ep.Path = RevocationPath

	data := url.Values{"token": []string{token}}
	hreq, err := http.NewRequest("POST", ep.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "")
	}

	hreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	hreq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	hres, err := r.Client.Do(hreq)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer hres.Body.Close()

	body, _ := ioutil.ReadAll(hres.Body)
	if hres.StatusCode < 200 || hres.StatusCode >= 300 {
		return errors.Errorf("Expected 2xx status code but got %d.\n%s", hres.StatusCode, body)
	}
	return nil
}
