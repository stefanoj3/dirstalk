package scan

import (
	"net/http"
	"net/url"

	"github.com/chuckpreslar/emission"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stefanoj3/dirstalk/pkg/dictionary"
	"github.com/stefanoj3/dirstalk/pkg/scan/client"
)

// StartScan is a convenience method that wires together all the dependencies needed to start a scan
func StartScan(logger *logrus.Logger, eventManager *emission.Emitter, cnf *Config, u *url.URL) error {
	c, err := client.NewClientFromConfig(
		cnf.TimeoutInMilliseconds,
		cnf.Socks5Url,
		cnf.UserAgent,
		cnf.UseCookieJar,
		cnf.Cookies,
		u,
	)
	if err != nil {
		return errors.Wrap(err, "failed to build client")
	}

	s := NewScanner(
		c,
		eventManager,
		logger,
	)

	dict, err := dictionary.NewDictionaryFrom(cnf.DictionaryPath, c)
	if err != nil {
		return errors.Wrap(err, "failed to build dictionary")
	}

	r := NewReProcessor(eventManager, cnf.HTTPMethods, dict)

	eventManager.On(EventResultFound, r.ReProcess)
	eventManager.On(EventTargetProduced, s.AddTarget)
	eventManager.On(EventProducerFinished, s.Release)

	targetProducer := NewTargetProducer(
		eventManager,
		cnf.HTTPMethods,
		dict,
		cnf.ScanDepth,
	)

	go targetProducer.Run()

	logger.WithFields(logrus.Fields{
		"url":               u.String(),
		"threads":           cnf.Threads,
		"dictionary.length": len(dict),
		"cookies":           strigifyCookies(cnf.Cookies),
		"user-agent":        cnf.UserAgent,
		"socks5":            cnf.Socks5Url,
	}).Info("Starting scan")

	s.Scan(u, cnf.Threads)

	logger.Info("Finished scan")

	return nil
}

func strigifyCookies(cookies []*http.Cookie) string {
	result := ""

	for _, cookie := range cookies {
		result += cookie.Name + "=" + cookie.Value + ";"
	}

	return result
}
