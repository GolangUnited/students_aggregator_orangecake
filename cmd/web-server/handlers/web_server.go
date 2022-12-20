package handlers

import (
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

type WebServer interface {
	Log() core.Logger
	DBReader() core.DBReader
	DBWriter() core.DBWriter
}

func NewWebServer(aLog core.Logger, aDBReader core.DBReader, aDBWriter core.DBWriter) (WebServer, error) {
	if aLog == nil {
		return nil, core.ErrLoggerNotAssigned
	}

	if aDBReader == nil {
		return nil, core.ErrDBReaderNotAssigned
	}

	if aDBWriter == nil {
		return nil, core.ErrDBWriterNotAssigned
	}

	lServer := serverImpl{
		log:      aLog,
		dbReader: aDBReader,
		dbWriter: aDBWriter,
	}

	return lServer, nil
}

type serverImpl struct {
	log      core.Logger
	dbReader core.DBReader
	dbWriter core.DBWriter
}

func (s serverImpl) Log() core.Logger {
	return s.log
}

func (s serverImpl) DBReader() core.DBReader {
	return s.dbReader
}

func (s serverImpl) DBWriter() core.DBWriter {
	return s.dbWriter
}
