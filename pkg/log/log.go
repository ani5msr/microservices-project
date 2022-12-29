package log

import (
	std_log "log"
	"os"

	kit_log "github.com/go-kit/log"
)

func NewLogger(service string) (logger kit_log.Logger) {
	w := kit_log.NewSyncWriter(os.Stderr)
	logger = kit_log.NewJSONLogger(w)
	logger = kit_log.With(logger, "service", service)
	logger = kit_log.With(logger, "timestamp", kit_log.DefaultTimestampUTC)
	logger = kit_log.With(logger, "called from", kit_log.DefaultCaller)

	return
}

func Fatal(v ...interface{}) {
	std_log.Fatal(v...)
}
