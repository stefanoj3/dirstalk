package scan

import (
	"net/http"
	"net/url"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/sirupsen/logrus"
)

func StartScan(logger *logrus.Logger, cnf *Config, u *url.URL) {
	eventManager := emission.NewEmitter()
	printer := NewResultLogger(logger)

	s := NewScanner(
		&http.Client{
			Timeout: time.Millisecond * time.Duration(cnf.TimeoutInMilliseconds),
		},
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
}
