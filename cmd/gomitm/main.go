package main

import (
	"flag"
	"github.com/qencept/gomitm/pkg/config"
	"github.com/qencept/gomitm/pkg/forgery"
	"github.com/qencept/gomitm/pkg/logger"
	"github.com/qencept/gomitm/pkg/proxy"
	"github.com/qencept/gomitm/pkg/session"
	"github.com/qencept/gomitm/pkg/session/copier"
	"github.com/qencept/gomitm/pkg/session/dump"
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

	//dohMutators := []doh.Mutator{NewApp(), dump3.NewCopy(l, cfg.Paths.Doh)}
	//http1Mutators := []http1.Mutator{doh.NewCopy(l, dohMutators...), dump2.NewCopy(l, cfg.Paths.Http)}
	//sessionMutators := []inspector.Mutator{http1.NewCopy(l, http1Mutators...), dump.NewCopy(l, cfg.Paths.Inspector)}
	sessionCreators := []session.Creator{copier.New(l), dump.New(l, cfg.Paths.Session), copier.New(l)}
	i := session.New(l, sessionCreators...)

	if err = proxy.New(a, t, f, i, l).Run(); err != nil {
		return err
	}

	return nil
}
