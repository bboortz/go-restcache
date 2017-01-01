package restcache

import (
	//	"github.com/davecgh/go-spew/spew"
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("restcache")

func LogAccess(method string, uri string, statusCode int, logTime time.Time) {
	route := getMethodName()
	LogAccessWithRoute(route, method, uri, statusCode, logTime)
}

func LogAccessWithRoute(route string, method string, uri string, statusCode int, logTime time.Time) {
	log.Infof("%s\t%s\t%s\t%s\t%d\t%s", log.Module, route, method, uri, statusCode, time.Since(logTime))
}
