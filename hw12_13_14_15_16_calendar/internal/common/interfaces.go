package common

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}
