package log

import (
	"github.com/sirupsen/logrus"
)

type Logger = logrus.Entry

type Factory struct {
	DefaultFields logrus.Fields
	Hooks         []logrus.Hook
	Logger        *logrus.Logger
}

func NewFactory(hooks ...logrus.Hook) *Factory {
	logger := logrus.New()
	for _, hook := range hooks {
		logger.Hooks.Add(hook)
	}
	return &Factory{Logger: logger, Hooks: hooks}
}

func (f *Factory) WithHooks(hooks ...logrus.Hook) *Factory {
	return NewFactory(append(f.Hooks, hooks...)...)
}

func (f *Factory) New(name string) *Logger {
	fields := logrus.Fields{"logger": name}
	for k, v := range f.DefaultFields {
		fields[k] = v
	}
	return logrus.New().WithFields(fields)
}
