package main

import (
	"flag"
	"github.com/qencept/gomitm/internal/config"
	"github.com/qencept/gomitm/internal/forging"
	"github.com/qencept/gomitm/internal/proxy"
	"github.com/qencept/gomitm/internal/trusted"
	"github.com/qencept/gomitm/mirror"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func main() {
	configPath := flag.String("config", "", "Configuration file")
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	cfg, err := config.ReadFile(*configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	trustedInstance, err := trusted.New(cfg.Proxy.TrustedRootCaCerts)
	if err != nil {
		logrus.Fatal(err)
	}

	forgingInstance, err := forging.New(cfg.Proxy.ForgingRootCa.Cert, cfg.Proxy.ForgingRootCa.Key)
	if err != nil {
		logrus.Fatal(err)
	}

	inspectors := []mirror.Inspector{mirror.NewDumper(), mirror.NewHttpParser()}
	shuttler := mirror.New(inspectors...)

	addr := ":" + cfg.Proxy.Port
	proxyInstance := proxy.New(addr, trustedInstance, forgingInstance, shuttler)
	if err := proxyInstance.Run(); err != nil {
		logrus.Fatal(err)
	}
}
