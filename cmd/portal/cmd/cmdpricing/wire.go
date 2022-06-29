//go:build wireinject
// +build wireinject

package cmdpricing

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/google/wire"

	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/infra/db"
	portalconfig "github.com/authgear/authgear-server/pkg/portal/config"
)

func NewStripeService(
	ctx context.Context,
	pool *db.Pool,
	databaseCredentials *config.DatabaseCredentials,
	stripeConfig *portalconfig.StripeConfig,
	hub *sentry.Hub,
) *StripeService {
	panic(wire.Build(DependencySet))
}
