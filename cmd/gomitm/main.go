package main

import (
	"flag"
	"github.com/qencept/gomitm/pkg/config"
	"github.com/qencept/gomitm/pkg/doh"
	ddump "github.com/qencept/gomitm/pkg/doh/dump"
	"github.com/qencept/gomitm/pkg/doh/tamper"
	"github.com/qencept/gomitm/pkg/forgery"
	"github.com/qencept/gomitm/pkg/http1"
	hdump "github.com/qencept/gomitm/pkg/http1/dump"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/proxy"
	"github.com/qencept/gomitm/pkg/session"
	sdump "github.com/qencept/gomitm/pkg/session/dump"
	"github.com/qencept/gomitm/pkg/trusted"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func main() {
	l := logrus.New()
	if err := run(l); err != nil {
		l.Fatal(err)
	}
}

func run(l logger.Logger) error {
	rand.Seed(time.Now().UnixNano())

	cfgFile := flag.String("config", "", "Configuration file")
	flag.Parse()
	cfg, err := config.ReadFile(*cfgFile)
	if err != nil {
		return err
	}

	a := ":" + cfg.Proxy.Port
	t, err := trusted.New(cfg.Proxy.TrustedRootCaCerts)
	if err != nil {
		return err
	}
	f, err := forgery.New(cfg.Proxy.ForgingRootCa.Cert, cfg.Proxy.ForgingRootCa.Key)
	if err != nil {
		return err
	}

	typeA := tamper.SubstitutionTypeA{"www.example.com.": [4]byte{1, 1, 1, 1}}
	dohCreators := []doh.Creator{ddump.New(l, cfg.Paths.Doh), tamper.New(typeA)}
	http1Creators := []http1.Creator{hdump.New(l, cfg.Paths.Http), doh.New(l, dohCreators...)}
	sessionCreators := []session.Creator{sdump.New(l, cfg.Paths.Session), http1.New(l, http1Creators...)}
	i := session.New(l, sessionCreators...)

	if err = proxy.New(a, t, f, i, l).Run(); err != nil {
		return err
	}

	return nil
}
