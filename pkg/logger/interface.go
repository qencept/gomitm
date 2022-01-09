package logger

type Logger interface {
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
}
