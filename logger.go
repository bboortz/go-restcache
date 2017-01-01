package restcache

import (
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("go-router")

func logAccess(route string, method string, uri string, statusCode int, logTime time.Time) {
	//	var log = logging.MustGetLogger("central-router")
	log.Infof("%s\t%s\t%s\t%d\t%s", route, method, uri, statusCode, time.Since(logTime))
}

func logServiceRegistered(s Service) {
	log.Infof("service registered: %d - %s ", s.Id, s.Name)
}
