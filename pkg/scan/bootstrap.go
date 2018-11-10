package scan

import (
	"net/http"
	"net/url"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

func StartScan(logger *logrus.Logger, cnf *Config, u *url.URL) error {
	eventManager := emission.NewEmitter()
	printer := NewResultLogger(logger)

	c, err := buildClientFrom(cnf)
	if err != nil {
		return errors.Wrap(err, "failed to build client")
	}

	s := NewScanner(
		c,
		eventManager,
		logger,
	)

	r := NewReProcessor(eventManager, cnf.HttpMethods, cnf.Dictionary)

	eventManager.On(EventResultFound, printer.Log)
	eventManager.On(EventResultFound, r.ReProcess)
	eventManager.On(EventTargetProduced, s.AddTarget)
	eventManager.On(EventProducerFinished, s.Release)

	targetProducer := NewTargetProducer(
		eventManager,
		cnf.HttpMethods,
		cnf.Dictionary,
		cnf.ScanDepth,
	)

	go targetProducer.Run()

	logger.WithFields(logrus.Fields{
		"url":               u.String(),
		"threads":           cnf.Threads,
		"dictionary.length": len(cnf.Dictionary),
	}).Info("Starting scan")

	s.Scan(u, cnf.Threads)

	logger.Info("Finished scan")

	return nil
}

func buildClientFrom(cnf *Config) (*http.Client, error) {
	c := &http.Client{
		Timeout: time.Millisecond * time.Duration(cnf.TimeoutInMilliseconds),
	}

	if cnf.Socks5Url != nil {
		tbDialer, err := proxy.FromURL(cnf.Socks5Url, proxy.Direct)
		if err != nil {
			return nil, err
		}

		tbTransport := &http.Transport{Dial: tbDialer.Dial}
		c.Transport = tbTransport
	}

	return c, nil
}
