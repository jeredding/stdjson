package stdjson

import (
	log "github.com/Sirupsen/logrus"
)

type LogFormatter struct {
	jf log.JSONFormatter
}

func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	entry.Time = entry.Time.UTC()
	return f.jf.Format(entry)
}
