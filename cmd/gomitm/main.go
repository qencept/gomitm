package main

import (
	"flag"
	"github.com/qencept/gomitm/internal/config"
	"github.com/qencept/gomitm/pkg/mirror/doh"
	"github.com/qencept/gomitm/pkg/mirror/http"
	"github.com/qencept/gomitm/pkg/mirror/session"
	"github.com/qencept/gomitm/pkg/mitm/forgery"
	"github.com/qencept/gomitm/pkg/mitm/proxy"
	"github.com/qencept/gomitm/pkg/mitm/shuttle"
	"github.com/qencept/gomitm/pkg/mitm/trusted"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
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

	addr := ":" + c.Proxy.Port
	if err := proxy.New(addr, t, f, setInspectors()).Run(); err != nil {
		return err
	}

	return nil
}

func setInspectors() shuttle.Shuttle {
	httpInspectors := []httpi.Inspector{httpi.NewDumper("http"), doh.NewDoh("doh")}
	sessionInspectors := []session.Inspector{session.NewDumper("session"), httpi.NewHttp(httpInspectors...)}
	return session.New(sessionInspectors...)
}
