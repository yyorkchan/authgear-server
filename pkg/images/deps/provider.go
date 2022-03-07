package deps

import (
	"net/http"

	getsentry "github.com/getsentry/sentry-go"

	imagesconfig "github.com/authgear/authgear-server/pkg/images/config"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/authgear/authgear-server/pkg/util/sentry"
)

type RootProvider struct {
	EnvironmentConfig imagesconfig.EnvironmentConfig
	ObjectStoreConfig *imagesconfig.ObjectStoreConfig
	LoggerFactory     *log.Factory
	SentryHub         *getsentry.Hub
}

func NewRootProvider(
	envConfig imagesconfig.EnvironmentConfig,
	objectStoreConfig *imagesconfig.ObjectStoreConfig,
) (*RootProvider, error) {
	logLevel, err := log.ParseLevel(string(envConfig.LogLevel))
	if err != nil {
		return nil, err
	}

	sentryHub, err := sentry.NewHub(string(envConfig.SentryDSN))
	if err != nil {
		return nil, err
	}

	loggerFactory := log.NewFactory(
		logLevel,
		log.NewDefaultMaskLogHook(),
		sentry.NewLogHookFromHub(sentryHub),
	)

	return &RootProvider{
		EnvironmentConfig: envConfig,
		ObjectStoreConfig: objectStoreConfig,
		LoggerFactory:     loggerFactory,
		SentryHub:         sentryHub,
	}, nil
}

type RequestProvider struct {
	RootProvider *RootProvider
	Request      *http.Request
}

func (p *RootProvider) Middleware(f func(*RequestProvider) httproute.Middleware) httproute.Middleware {
	return httproute.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestProvider := &RequestProvider{
				RootProvider: p,
				Request:      r,
			}
			m := f(requestProvider)
			h := m.Handle(next)
			h.ServeHTTP(w, r)
		})
	})
}

func (p *RootProvider) Handler(f func(*RequestProvider) http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestProvider := &RequestProvider{
			RootProvider: p,
			Request:      r,
		}
		h := f(requestProvider)
		h.ServeHTTP(w, r)
	})
}
