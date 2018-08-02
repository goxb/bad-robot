package log

import (
	"fmt"
	"path"

	"github.com/cihub/seelog"
	"github.com/smallnest/rpcx/log"
)

type rpcxLogger struct {
	logger seelog.LoggerInterface
}

func (l *rpcxLogger) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}

func (l *rpcxLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

func (l *rpcxLogger) Info(v ...interface{}) {
	l.logger.Info(v...)
}

func (l *rpcxLogger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

func (l *rpcxLogger) Warn(v ...interface{}) {
	l.logger.Warn(v...)
}

func (l *rpcxLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

func (l *rpcxLogger) Error(v ...interface{}) {
	l.logger.Error(v...)
}

func (l *rpcxLogger) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

func (l *rpcxLogger) Fatal(v ...interface{}) {
	l.logger.Error(v...)
}

func (l *rpcxLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

func (l *rpcxLogger) Panic(v ...interface{}) {
	l.logger.Critical(v...)
}

func (l *rpcxLogger) Panicf(format string, v ...interface{}) {
	l.logger.Criticalf(format, v...)
}

func InitLog(loglevel string, logfile string) int {
	tostdflag := false
	logConf := `
<seelog minlevel="` + loglevel + `">
    <outputs formatid="main">
        <filter levels="critical">
            <file path="` + path.Dir(logfile) + `/stack.log"/>
        </filter>
        <filter levels="trace, debug, info, warn, error">
`
	if tostdflag {
		logConf = logConf + "<console />"
	}
	logConf = logConf + `
            <rollingfile type="date" filename="` + logfile + `" datepattern="02.01.2006" maxrolls="1"/>
        </filter>
    </outputs>
    <formats>
        <format id="main" format="%Date(2006-01-02T15:04:05.999999999Z07:00) [@%File.%Line][%LEV] %Msg%n"/>
    </formats>
</seelog>
`
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(logConf))
	if err != nil {
		fmt.Println("err parsing config:", err)
		return -1
	}

	log.SetLogger(&rpcxLogger{
		logger: logger,
	})

	seelog.ReplaceLogger(logger)
	return 0
}
