/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/x"
	"github.com/ory/x/flagx"
	"github.com/ory/x/logrusx"

	"github.com/julienschmidt/httprouter"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"

	"github.com/ory/graceful"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/metricsx"
)

var _ = &consent.Handler{}

func EnhanceMiddleware(d driver.Driver, n *negroni.Negroni, address string, router *httprouter.Router, enableCORS bool, iface string) http.Handler {
	if !x.AddressIsUnixSocket(address) {
		n.UseFunc(x.RejectInsecureRequests(d.Registry(), d.Configuration()))
	}
	n.UseHandler(router)
	if enableCORS {
		options := d.Configuration().CORSOptions(iface)
		d.Registry().Logger().
			WithField("options", fmt.Sprintf("%+v", options)).
			Infof("Enabling CORS on interface: %s", address)
		return cors.New(options).Handler(n)
	}
	return n
}

func isDSNAllowed(d driver.Driver) {
	if d.Configuration().DSN() == "memory" {
		d.Registry().Logger().Fatalf(`When using "hydra serve admin" or "hydra serve public" the DSN can not be set to "memory".`)
	}
}

func RunServeAdmin(version, build, date string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		d := driver.NewDefaultDriver(
			logrusx.New(),
			prometheus.NewRegistry(),
			flagx.MustGetBool(cmd, "dangerous-force-http"),
			flagx.MustGetStringSlice(cmd, "dangerous-allow-insecure-redirect-urls"),
			version, build, date, true,
		).CallRegistry()

		isDSNAllowed(d)

		admin, _, adminmw, _ := setup(d, cmd)
		cert := getOrCreateTLSCertificate(cmd, d) // we do not want to run this concurrently.

		var wg sync.WaitGroup
		wg.Add(1)

		go serve(d, cmd, &wg,
			EnhanceMiddleware(d, adminmw, d.Configuration().AdminListenOn(), admin.Router, d.Configuration().CORSEnabled("admin"), "admin"),
			d.Configuration().AdminListenOn(), cert,
		)

		wg.Wait()
	}
}

func RunServePublic(version, build, date string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		d := driver.NewDefaultDriver(
			logrusx.New(),
			prometheus.NewRegistry(),
			flagx.MustGetBool(cmd, "dangerous-force-http"),
			flagx.MustGetStringSlice(cmd, "dangerous-allow-insecure-redirect-urls"),
			version, build, date, true,
		).CallRegistry()

		isDSNAllowed(d)

		_, public, _, publicmw := setup(d, cmd)
		cert := getOrCreateTLSCertificate(cmd, d) // we do not want to run this concurrently.

		var wg sync.WaitGroup
		wg.Add(1)

		go serve(d, cmd, &wg,
			EnhanceMiddleware(d, publicmw, d.Configuration().PublicListenOn(), public.Router, false, "public"),
			d.Configuration().PublicListenOn(), cert,
		)

		wg.Wait()
	}
}

func RunServeAll(version, build, date string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		d := driver.NewDefaultDriver(
			logrusx.New(),
			prometheus.NewRegistry(),
			flagx.MustGetBool(cmd, "dangerous-force-http"),
			flagx.MustGetStringSlice(cmd, "dangerous-allow-insecure-redirect-urls"),
			version, build, date, true,
		).CallRegistry()

		admin, public, adminmw, publicmw := setup(d, cmd)
		cert := getOrCreateTLSCertificate(cmd, d) // we do not want to run this concurrently.

		var wg sync.WaitGroup
		wg.Add(2)

		go serve(d, cmd, &wg,
			EnhanceMiddleware(d, publicmw, d.Configuration().PublicListenOn(), public.Router, false, "public"),
			d.Configuration().PublicListenOn(), cert,
		)

		go serve(d, cmd, &wg,
			EnhanceMiddleware(d, adminmw, d.Configuration().AdminListenOn(), admin.Router, d.Configuration().CORSEnabled("admin"), "admin"),
			d.Configuration().AdminListenOn(), cert,
		)

		wg.Wait()
	}
}

func setup(d driver.Driver, cmd *cobra.Command) (admin *x.RouterAdmin, public *x.RouterPublic, adminmw, publicmw *negroni.Negroni) {
	fmt.Println(banner(d.Registry().BuildVersion()))

	adminmw = negroni.New()
	publicmw = negroni.New()

	admin = x.NewRouterAdmin()
	public = x.NewRouterPublic()

	if tracer := d.Registry().Tracer(); tracer.IsLoaded() {
		adminmw.Use(tracer)
		publicmw.Use(tracer)
	}

	adminmw.Use(negronilogrus.NewMiddlewareFromLogger(d.Registry().Logger().(*logrus.Logger), fmt.Sprintf("hydra/admin: %s", d.Configuration().IssuerURL().String())))
	adminmw.Use(d.Registry().PrometheusManager())

	publicmw.Use(negronilogrus.NewMiddlewareFromLogger(d.Registry().Logger().(*logrus.Logger), fmt.Sprintf("hydra/public: %s", d.Configuration().IssuerURL().String())))
	publicmw.Use(d.Registry().PrometheusManager())

	var (
		bridges = []metricsx.Bridge{}
		l       = d.Registry().Logger()
		ctx     = context.Background()
	)

	optOut, err := cmd.Flags().GetBool("sqa-opt-out")
	if !optOut {
		optOut, err = cmd.Flags().GetBool("disable-telemetry")
		if optOut {
			l.Warn(`Command line argument "--disable-telemetry" has been deprecated and will be removed in an upcoming release. Use "--sqa-opt-out" instead.`)
		}
	}
	cmdx.Must(err, "Unable to get command line flag.")

	if !optOut {
		optOut = viper.GetBool("sqa.opt_out")
	}

	if !optOut {
		optOut = viper.GetBool("DISABLE_TELEMETRY")
		if optOut {
			l.Warn(`Environment variable "DISABLE_TELEMETRY" has been deprecated and will be removed in an upcoming release. Use configuration key "sqa.opt_out: true" or environment variable "SQA_OPT_OUT=true" instead.`)
		}
	}

	if !optOut {
		l.Info("Software quality assurance features are enabled. Learn more at: https://www.ory.sh/docs/ecosystem/sqa")
		s := &metricsx.SegmentOptions{
			Service: "ory-hydra",
			ClusterID: metricsx.Hash(fmt.Sprintf("%s|%s",
				d.Configuration().IssuerURL().String(),
				d.Configuration().DSN(),
			)),
			IsDevelopment: d.Configuration().DSN() == "memory" ||
				d.Configuration().IssuerURL().String() == "" ||
				strings.Contains(d.Configuration().IssuerURL().String(), "localhost"),
			WhitelistedPaths: []string{
				jwk.KeyHandlerPath,
				jwk.WellKnownKeysPath,

				client.ClientsHandlerPath,

				oauth2.DefaultConsentPath,
				oauth2.DefaultLoginPath,
				oauth2.DefaultPostLogoutPath,
				oauth2.DefaultLogoutPath,
				oauth2.DefaultErrorPath,
				oauth2.TokenPath,
				oauth2.AuthPath,
				oauth2.LogoutPath,
				oauth2.UserinfoPath,
				oauth2.WellKnownPath,
				oauth2.JWKPath,
				oauth2.IntrospectPath,
				oauth2.RevocationPath,
				oauth2.FlushPath,

				consent.ConsentPath,
				consent.ConsentPath + "/accept",
				consent.ConsentPath + "/reject",
				consent.LoginPath,
				consent.LoginPath + "/accept",
				consent.LoginPath + "/reject",
				consent.LogoutPath,
				consent.LogoutPath + "/accept",
				consent.LogoutPath + "/reject",
				consent.SessionsPath + "/login",
				consent.SessionsPath + "/consent",

				healthx.AliveCheckPath,
				healthx.ReadyCheckPath,
				healthx.VersionPath,
				driver.MetricsPrometheusPath,
				"/",
			},
			WriteKey:     "h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr",
			BuildVersion: d.Registry().BuildVersion(),
			BuildTime:    d.Registry().BuildDate(),
			BuildHash:    d.Registry().BuildHash(),
		}
		b, err := metricsx.NewFormattedSegmentBridge(ctx, s, l, d.Registry().PrometheusManager().Registry)
		if err != nil {
			l.WithError(err).Debug("Software quality assurance features could not be started. Learn more at: https://www.ory.sh/docs/ecosystem/sqa")
		}

		publicmw.Use(b)
		bridges = append(bridges, b)
	}
	bm := metricsx.NewBridgeManager(&metricsx.BridgeManagerOptions{
		Interval: time.Duration(5 * time.Second),
		Log:      l,
	}, bridges)

	go bm.Start(ctx)

	// metrics := metricsx.New(
	// 	cmd,
	// 	d.Registry().Logger(),
	// 	&metricsx.Options{
	// 		Service: "ory-hydra",
	// 		ClusterID: metricsx.Hash(fmt.Sprintf("%s|%s",
	// 			d.Configuration().IssuerURL().String(),
	// 			d.Configuration().DSN(),
	// 		)),
	// 		IsDevelopment: d.Configuration().DSN() == "memory" ||
	// 			d.Configuration().IssuerURL().String() == "" ||
	// 			strings.Contains(d.Configuration().IssuerURL().String(), "localhost"),
	// 		WriteKey: "h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr",
	// 		WhitelistedPaths: []string{
	// 			client.ClientsHandlerPath,
	// 			jwk.KeyHandlerPath,
	// 			jwk.WellKnownKeysPath,
	// 			oauth2.DefaultConsentPath,
	// 			oauth2.TokenPath,
	// 			oauth2.AuthPath,
	// 			oauth2.UserinfoPath,
	// 			oauth2.WellKnownPath,
	// 			oauth2.IntrospectPath,
	// 			oauth2.RevocationPath,
	// 			consent.ConsentPath,
	// 			consent.LoginPath,
	// 			healthx.AliveCheckPath,
	// 			healthx.ReadyCheckPath,
	// 			healthx.VersionPath,
	// 			driver.MetricsPrometheusPath,
	// 			"/oauth2/auth/sessions/login",
	// 			"/oauth2/auth/sessions/consent",
	// 			"/",
	// 		},
	// 		BuildVersion: d.Registry().BuildVersion(),
	// 		BuildTime:    d.Registry().BuildDate(),
	// 		BuildHash:    d.Registry().BuildHash(),
	// 	},
	// )

	// adminmw.Use(metrics)
	// publicmw.Use(metrics)

	d.Registry().RegisterRoutes(admin, public)

	return
}

func serve(d driver.Driver, cmd *cobra.Command, wg *sync.WaitGroup, handler http.Handler, address string, cert []tls.Certificate) {
	defer wg.Done()

	var srv = graceful.WithDefaults(&http.Server{
		Addr:    address,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: cert,
		},
	})

	if d.Registry().Tracer().IsLoaded() {
		srv.RegisterOnShutdown(d.Registry().Tracer().Close)
	}

	if err := graceful.Graceful(func() error {
		var err error
		d.Registry().Logger().Infof("Setting up http server on %s", address)
		if x.AddressIsUnixSocket(address) {
			addr := strings.TrimPrefix(address, "unix:")
			unixListener, e := net.Listen("unix", addr)
			if e != nil {
				return e
			}
			err = srv.Serve(unixListener)
		} else {
			if !d.Configuration().ServesHTTPS() {
				d.Registry().Logger().Warnln("HTTPS disabled. Never do this in production.")
				err = srv.ListenAndServe()
			} else if len(d.Configuration().AllowTLSTerminationFrom()) > 0 {
				d.Registry().Logger().Infoln("TLS termination enabled, disabling https.")
				err = srv.ListenAndServe()
			} else {
				err = srv.ListenAndServeTLS("", "")
			}
		}

		return err
	}, srv.Shutdown); err != nil {
		d.Registry().Logger().WithError(err).Fatal("Could not gracefully run server")
	}
}
