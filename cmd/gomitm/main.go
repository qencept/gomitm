package main

import (
	"flag"
	"github.com/qencept/gomitm/internal/config"
	"github.com/qencept/gomitm/pkg/mirror/doh"
	"github.com/qencept/gomitm/pkg/mirror/http1"
	"github.com/qencept/gomitm/pkg/mirror/http1dump"
	"github.com/qencept/gomitm/pkg/mirror/session"
	"github.com/qencept/gomitm/pkg/mirror/sessiondump"
	"github.com/qencept/gomitm/pkg/mitm/forgery"
	"github.com/qencept/gomitm/pkg/mitm/proxy"
	"github.com/qencept/gomitm/pkg/mitm/trusted"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	if err := run(logger); err != nil {
		logger.Fatal(err)
	}
}

func run(l *logrus.Logger) error {
	rand.Seed(time.Now().UnixNano())

	cfgFile := flag.String("config", "", "Configuration file")
	flag.Parse()

	c, err := config.ReadFile(*cfgFile)
	if err != nil {
		return err
	}

	t, err := trusted.New(c.Proxy.TrustedRootCaCerts)
	if err != nil {
		return err
	}

	f, err := forgery.New(c.Proxy.ForgingRootCa.Cert, c.Proxy.ForgingRootCa.Key)
	if err != nil {
		return err
	}

	httpInspectors := []http1.Inspector{http1dump.New(c.Paths.Http), doh.New(c.Paths.Doh)}
	sessionInspectors := []session.Inspector{sessiondump.New(c.Paths.Session), http1.New(l, httpInspectors...)}
	s := session.New(l, sessionInspectors...)

	a := ":" + c.Proxy.Port
	if err = proxy.New(a, t, f, l, s).Run(); err != nil {
		return err
	}

	return nil
}
