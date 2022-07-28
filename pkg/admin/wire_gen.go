// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package admin

import (
	"context"
	facade2 "github.com/authgear/authgear-server/pkg/admin/facade"
	"github.com/authgear/authgear-server/pkg/admin/graphql"
	"github.com/authgear/authgear-server/pkg/admin/loader"
	service3 "github.com/authgear/authgear-server/pkg/admin/service"
	"github.com/authgear/authgear-server/pkg/admin/transport"
	"github.com/authgear/authgear-server/pkg/lib/admin/authz"
	"github.com/authgear/authgear-server/pkg/lib/audit"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticationinfo"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/oob"
	passkey2 "github.com/authgear/authgear-server/pkg/lib/authn/authenticator/passkey"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/password"
	service2 "github.com/authgear/authgear-server/pkg/lib/authn/authenticator/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/totp"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/whatsapp"
	"github.com/authgear/authgear-server/pkg/lib/authn/challenge"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/anonymous"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/biometric"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/loginid"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/oauth"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/passkey"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/mfa"
	"github.com/authgear/authgear-server/pkg/lib/authn/otp"
	"github.com/authgear/authgear-server/pkg/lib/authn/sso"
	stdattrs2 "github.com/authgear/authgear-server/pkg/lib/authn/stdattrs"
	"github.com/authgear/authgear-server/pkg/lib/authn/user"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/elasticsearch"
	"github.com/authgear/authgear-server/pkg/lib/event"
	"github.com/authgear/authgear-server/pkg/lib/facade"
	"github.com/authgear/authgear-server/pkg/lib/feature/customattrs"
	"github.com/authgear/authgear-server/pkg/lib/feature/forgotpassword"
	"github.com/authgear/authgear-server/pkg/lib/feature/stdattrs"
	"github.com/authgear/authgear-server/pkg/lib/feature/verification"
	"github.com/authgear/authgear-server/pkg/lib/feature/welcomemessage"
	"github.com/authgear/authgear-server/pkg/lib/healthz"
	"github.com/authgear/authgear-server/pkg/lib/hook"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/appdb"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/auditdb"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/globaldb"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/lib/interaction"
	"github.com/authgear/authgear-server/pkg/lib/nonce"
	oauth2 "github.com/authgear/authgear-server/pkg/lib/oauth"
	"github.com/authgear/authgear-server/pkg/lib/oauth/pq"
	"github.com/authgear/authgear-server/pkg/lib/oauth/redis"
	"github.com/authgear/authgear-server/pkg/lib/presign"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/lib/session"
	"github.com/authgear/authgear-server/pkg/lib/session/access"
	"github.com/authgear/authgear-server/pkg/lib/session/idpsession"
	"github.com/authgear/authgear-server/pkg/lib/translation"
	"github.com/authgear/authgear-server/pkg/lib/tutorial"
	"github.com/authgear/authgear-server/pkg/lib/usage"
	"github.com/authgear/authgear-server/pkg/lib/web"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/authgear/authgear-server/pkg/util/rand"
	"github.com/authgear/authgear-server/pkg/util/template"
	"net/http"
)

// Injectors from wire.go:

func newPanicMiddleware(p *deps.RootProvider) httproute.Middleware {
	factory := p.LoggerFactory
	panicMiddlewareLogger := middleware.NewPanicMiddlewareLogger(factory)
	panicMiddleware := &middleware.PanicMiddleware{
		Logger: panicMiddlewareLogger,
	}
	return panicMiddleware
}

func newHealthzHandler(p *deps.RootProvider, w http.ResponseWriter, r *http.Request, ctx context.Context) http.Handler {
	pool := p.DatabasePool
	environmentConfig := p.EnvironmentConfig
	globalDatabaseCredentialsEnvironmentConfig := &environmentConfig.GlobalDatabase
	databaseEnvironmentConfig := &environmentConfig.DatabaseConfig
	factory := p.LoggerFactory
	handle := globaldb.NewHandle(ctx, pool, globalDatabaseCredentialsEnvironmentConfig, databaseEnvironmentConfig, factory)
	sqlExecutor := globaldb.NewSQLExecutor(ctx, handle)
	handlerLogger := healthz.NewHandlerLogger(factory)
	handler := &healthz.Handler{
		Context:        ctx,
		GlobalDatabase: handle,
		GlobalExecutor: sqlExecutor,
		Logger:         handlerLogger,
	}
	return handler
}

func newSentryMiddleware(p *deps.RootProvider) httproute.Middleware {
	hub := p.SentryHub
	environmentConfig := p.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	sentryMiddleware := &middleware.SentryMiddleware{
		SentryHub:  hub,
		TrustProxy: trustProxy,
	}
	return sentryMiddleware
}

func newBodyLimitMiddleware(p *deps.RootProvider) httproute.Middleware {
	bodyLimitMiddleware := &middleware.BodyLimitMiddleware{}
	return bodyLimitMiddleware
}

func newAuthorizationMiddleware(p *deps.RequestProvider, auth config.AdminAPIAuth) httproute.Middleware {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	logger := authz.NewLogger(factory)
	configConfig := appProvider.Config
	appConfig := configConfig.AppConfig
	appID := appConfig.ID
	secretConfig := configConfig.SecretConfig
	adminAPIAuthKey := deps.ProvideAdminAPIAuthKeyMaterials(secretConfig)
	clock := _wireSystemClockValue
	authzMiddleware := &authz.Middleware{
		Logger:  logger,
		Auth:    auth,
		AppID:   appID,
		AuthKey: adminAPIAuthKey,
		Clock:   clock,
	}
	return authzMiddleware
}

var (
	_wireSystemClockValue = clock.NewSystemClock()
)

func newGraphQLHandler(p *deps.RequestProvider) http.Handler {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	logger := graphql.NewLogger(factory)
	configConfig := appProvider.Config
	secretConfig := configConfig.SecretConfig
	databaseCredentials := deps.ProvideDatabaseCredentials(secretConfig)
	appConfig := configConfig.AppConfig
	appID := appConfig.ID
	sqlBuilderApp := appdb.NewSQLBuilderApp(databaseCredentials, appID)
	request := p.Request
	contextContext := deps.ProvideRequestContext(request)
	handle := appProvider.AppDatabase
	sqlExecutor := appdb.NewSQLExecutor(contextContext, handle)
	clockClock := _wireSystemClockValue
	store := &user.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
		Clock:       clockClock,
	}
	rawQueries := &user.RawQueries{
		Store: store,
	}
	authenticationConfig := appConfig.Authentication
	identityConfig := appConfig.Identity
	featureConfig := configConfig.FeatureConfig
	identityFeatureConfig := featureConfig.Identity
	serviceStore := &service.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	loginidStore := &loginid.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	loginIDConfig := identityConfig.LoginID
	manager := appProvider.Resources
	typeCheckerFactory := &loginid.TypeCheckerFactory{
		Config:    loginIDConfig,
		Resources: manager,
	}
	checker := &loginid.Checker{
		Config:             loginIDConfig,
		TypeCheckerFactory: typeCheckerFactory,
	}
	normalizerFactory := &loginid.NormalizerFactory{
		Config: loginIDConfig,
	}
	provider := &loginid.Provider{
		Store:             loginidStore,
		Config:            loginIDConfig,
		Checker:           checker,
		NormalizerFactory: normalizerFactory,
		Clock:             clockClock,
	}
	oauthStore := &oauth.Store{
		SQLBuilder:     sqlBuilderApp,
		SQLExecutor:    sqlExecutor,
		IdentityConfig: identityConfig,
	}
	oauthProvider := &oauth.Provider{
		Store:          oauthStore,
		Clock:          clockClock,
		IdentityConfig: identityConfig,
	}
	anonymousStore := &anonymous.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	anonymousProvider := &anonymous.Provider{
		Store: anonymousStore,
		Clock: clockClock,
	}
	biometricStore := &biometric.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	biometricProvider := &biometric.Provider{
		Store: biometricStore,
		Clock: clockClock,
	}
	passkeyStore := &passkey.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	passkeyProvider := &passkey.Provider{
		Store: passkeyStore,
		Clock: clockClock,
	}
	serviceService := &service.Service{
		Authentication:        authenticationConfig,
		Identity:              identityConfig,
		IdentityFeatureConfig: identityFeatureConfig,
		Store:                 serviceStore,
		LoginID:               provider,
		OAuth:                 oauthProvider,
		Anonymous:             anonymousProvider,
		Biometric:             biometricProvider,
		Passkey:               passkeyProvider,
	}
	store2 := &service2.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	passwordStore := &password.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	authenticatorConfig := appConfig.Authenticator
	authenticatorPasswordConfig := authenticatorConfig.Password
	passwordLogger := password.NewLogger(factory)
	historyStore := &password.HistoryStore{
		Clock:       clockClock,
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	authenticatorFeatureConfig := featureConfig.Authenticator
	passwordChecker := password.ProvideChecker(authenticatorPasswordConfig, authenticatorFeatureConfig, historyStore)
	housekeeperLogger := password.NewHousekeeperLogger(factory)
	housekeeper := &password.Housekeeper{
		Store:  historyStore,
		Logger: housekeeperLogger,
		Config: authenticatorPasswordConfig,
	}
	passwordProvider := &password.Provider{
		Store:           passwordStore,
		Config:          authenticatorPasswordConfig,
		Clock:           clockClock,
		Logger:          passwordLogger,
		PasswordHistory: historyStore,
		PasswordChecker: passwordChecker,
		Housekeeper:     housekeeper,
	}
	store3 := &passkey2.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	provider2 := &passkey2.Provider{
		Store: store3,
		Clock: clockClock,
	}
	totpStore := &totp.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	authenticatorTOTPConfig := authenticatorConfig.TOTP
	totpProvider := &totp.Provider{
		Store:  totpStore,
		Config: authenticatorTOTPConfig,
		Clock:  clockClock,
	}
	authenticatorOOBConfig := authenticatorConfig.OOB
	oobStore := &oob.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	appredisHandle := appProvider.Redis
	storeRedis := &oob.StoreRedis{
		Redis: appredisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	oobLogger := oob.NewLogger(factory)
	oobProvider := &oob.Provider{
		Config:    authenticatorOOBConfig,
		Store:     oobStore,
		CodeStore: storeRedis,
		Clock:     clockClock,
		Logger:    oobLogger,
	}
	ratelimitLogger := ratelimit.NewLogger(factory)
	storageRedis := &ratelimit.StorageRedis{
		AppID: appID,
		Redis: appredisHandle,
	}
	limiter := &ratelimit.Limiter{
		Logger:  ratelimitLogger,
		Storage: storageRedis,
		Clock:   clockClock,
	}
	service4 := &service2.Service{
		Store:       store2,
		Password:    passwordProvider,
		Passkey:     provider2,
		TOTP:        totpProvider,
		OOBOTP:      oobProvider,
		RateLimiter: limiter,
	}
	rootProvider := appProvider.RootProvider
	environmentConfig := rootProvider.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	remoteIP := deps.ProvideRemoteIP(request, trustProxy)
	verificationLogger := verification.NewLogger(factory)
	verificationConfig := appConfig.Verification
	userProfileConfig := appConfig.UserProfile
	verificationStoreRedis := &verification.StoreRedis{
		Redis: appredisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	storePQ := &verification.StorePQ{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	verificationService := &verification.Service{
		RemoteIP:          remoteIP,
		Logger:            verificationLogger,
		Config:            verificationConfig,
		UserProfileConfig: userProfileConfig,
		Clock:             clockClock,
		CodeStore:         verificationStoreRedis,
		ClaimStore:        storePQ,
		RateLimiter:       limiter,
	}
	httpProto := deps.ProvideHTTPProto(request, trustProxy)
	httpHost := deps.ProvideHTTPHost(request, trustProxy)
	imagesCDNHost := environmentConfig.ImagesCDNHost
	pictureTransformer := &stdattrs.PictureTransformer{
		HTTPProto:     httpProto,
		HTTPHost:      httpHost,
		ImagesCDNHost: imagesCDNHost,
	}
	serviceNoEvent := &stdattrs.ServiceNoEvent{
		UserProfileConfig: userProfileConfig,
		Identities:        serviceService,
		UserQueries:       rawQueries,
		UserStore:         store,
		ClaimStore:        storePQ,
		Transformer:       pictureTransformer,
	}
	customattrsServiceNoEvent := &customattrs.ServiceNoEvent{
		Config:      userProfileConfig,
		UserQueries: rawQueries,
		UserStore:   store,
	}
	queries := &user.Queries{
		RawQueries:         rawQueries,
		Store:              store,
		Identities:         serviceService,
		Authenticators:     service4,
		Verification:       verificationService,
		StandardAttributes: serviceNoEvent,
		CustomAttributes:   customattrsServiceNoEvent,
	}
	userLoader := loader.NewUserLoader(queries)
	identityLoader := loader.NewIdentityLoader(serviceService)
	authenticatorLoader := loader.NewAuthenticatorLoader(service4)
	readHandle := appProvider.AuditReadDatabase
	auditDatabaseCredentials := deps.ProvideAuditDatabaseCredentials(secretConfig)
	auditdbSQLBuilderApp := auditdb.NewSQLBuilderApp(auditDatabaseCredentials, appID)
	readSQLExecutor := auditdb.NewReadSQLExecutor(contextContext, readHandle)
	readStore := &audit.ReadStore{
		SQLBuilder:  auditdbSQLBuilderApp,
		SQLExecutor: readSQLExecutor,
	}
	query := &audit.Query{
		Database: readHandle,
		Store:    readStore,
	}
	auditLogLoader := loader.NewAuditLogLoader(query)
	elasticsearchCredentials := deps.ProvideElasticsearchCredentials(secretConfig)
	client := elasticsearch.NewClient(elasticsearchCredentials)
	queue := appProvider.TaskQueue
	elasticsearchService := &elasticsearch.Service{
		AppID:     appID,
		Client:    client,
		Users:     queries,
		OAuth:     oauthStore,
		LoginID:   loginidStore,
		TaskQueue: queue,
	}
	defaultLanguageTag := deps.ProvideDefaultLanguageTag(configConfig)
	supportedLanguageTags := deps.ProvideSupportedLanguageTags(configConfig)
	resolver := &template.Resolver{
		Resources:             manager,
		DefaultLanguageTag:    defaultLanguageTag,
		SupportedLanguageTags: supportedLanguageTags,
	}
	engine := &template.Engine{
		Resolver: resolver,
	}
	httpConfig := appConfig.HTTP
	localizationConfig := appConfig.Localization
	staticAssetURLPrefix := environmentConfig.StaticAssetURLPrefix
	globalEmbeddedResourceManager := rootProvider.EmbeddedResources
	staticAssetResolver := &web.StaticAssetResolver{
		Context:            contextContext,
		Config:             httpConfig,
		Localization:       localizationConfig,
		StaticAssetsPrefix: staticAssetURLPrefix,
		Resources:          manager,
		EmbeddedResources:  globalEmbeddedResourceManager,
	}
	translationService := &translation.Service{
		Context:        contextContext,
		TemplateEngine: engine,
		StaticAssets:   staticAssetResolver,
	}
	welcomeMessageConfig := appConfig.WelcomeMessage
	userAgentString := deps.ProvideUserAgentString(request)
	eventLogger := event.NewLogger(factory)
	sqlBuilder := appdb.NewSQLBuilder(databaseCredentials)
	storeImpl := &event.StoreImpl{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	resolverImpl := &event.ResolverImpl{
		Users: queries,
	}
	hookLogger := hook.NewLogger(factory)
	hookConfig := appConfig.Hook
	webhookKeyMaterials := deps.ProvideWebhookKeyMaterials(secretConfig)
	syncHTTPClient := hook.NewSyncHTTPClient(hookConfig)
	asyncHTTPClient := hook.NewAsyncHTTPClient()
	deliverer := &hook.Deliverer{
		Config:             hookConfig,
		Secret:             webhookKeyMaterials,
		Clock:              clockClock,
		SyncHTTP:           syncHTTPClient,
		AsyncHTTP:          asyncHTTPClient,
		StandardAttributes: serviceNoEvent,
		CustomAttributes:   customattrsServiceNoEvent,
	}
	sink := &hook.Sink{
		Logger:    hookLogger,
		Deliverer: deliverer,
	}
	auditLogger := audit.NewLogger(factory)
	writeHandle := appProvider.AuditWriteDatabase
	writeSQLExecutor := auditdb.NewWriteSQLExecutor(contextContext, writeHandle)
	writeStore := &audit.WriteStore{
		SQLBuilder:  auditdbSQLBuilderApp,
		SQLExecutor: writeSQLExecutor,
	}
	auditSink := &audit.Sink{
		Logger:   auditLogger,
		Database: writeHandle,
		Store:    writeStore,
	}
	globalDatabaseCredentialsEnvironmentConfig := &environmentConfig.GlobalDatabase
	globaldbSQLBuilder := globaldb.NewSQLBuilder(globalDatabaseCredentialsEnvironmentConfig)
	pool := rootProvider.DatabasePool
	databaseEnvironmentConfig := &environmentConfig.DatabaseConfig
	globaldbHandle := globaldb.NewHandle(contextContext, pool, globalDatabaseCredentialsEnvironmentConfig, databaseEnvironmentConfig, factory)
	globaldbSQLExecutor := globaldb.NewSQLExecutor(contextContext, globaldbHandle)
	tutorialStoreImpl := &tutorial.StoreImpl{
		SQLBuilder:  globaldbSQLBuilder,
		SQLExecutor: globaldbSQLExecutor,
	}
	tutorialService := &tutorial.Service{
		Store: tutorialStoreImpl,
	}
	tutorialSink := &tutorial.Sink{
		AppID:        appID,
		Service:      tutorialService,
		GlobalHandle: globaldbHandle,
	}
	elasticsearchLogger := elasticsearch.NewLogger(factory)
	service5 := elasticsearch.Service{
		AppID:     appID,
		Client:    client,
		Users:     queries,
		OAuth:     oauthStore,
		LoginID:   loginidStore,
		TaskQueue: queue,
	}
	elasticsearchSink := &elasticsearch.Sink{
		Logger:   elasticsearchLogger,
		Service:  service5,
		Database: handle,
	}
	eventService := event.NewService(contextContext, remoteIP, userAgentString, eventLogger, handle, clockClock, localizationConfig, storeImpl, resolverImpl, sink, auditSink, tutorialSink, elasticsearchSink)
	welcomemessageProvider := &welcomemessage.Provider{
		Translation:          translationService,
		RateLimiter:          limiter,
		WelcomeMessageConfig: welcomeMessageConfig,
		TaskQueue:            queue,
		Events:               eventService,
	}
	rawCommands := &user.RawCommands{
		Store:                  store,
		Clock:                  clockClock,
		WelcomeMessageProvider: welcomemessageProvider,
	}
	commands := &user.Commands{
		RawCommands:        rawCommands,
		RawQueries:         rawQueries,
		Events:             eventService,
		Verification:       verificationService,
		UserProfileConfig:  userProfileConfig,
		StandardAttributes: serviceNoEvent,
		CustomAttributes:   customattrsServiceNoEvent,
	}
	userProvider := &user.Provider{
		Commands: commands,
		Queries:  queries,
	}
	storeDeviceTokenRedis := &mfa.StoreDeviceTokenRedis{
		Redis: appredisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	storeRecoveryCodePQ := &mfa.StoreRecoveryCodePQ{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	mfaService := &mfa.Service{
		DeviceTokens:  storeDeviceTokenRedis,
		RecoveryCodes: storeRecoveryCodePQ,
		Clock:         clockClock,
		Config:        authenticationConfig,
		RateLimiter:   limiter,
	}
	stdattrsService := &stdattrs.Service{
		UserProfileConfig: userProfileConfig,
		ServiceNoEvent:    serviceNoEvent,
		Identities:        serviceService,
		UserQueries:       rawQueries,
		UserStore:         store,
		Events:            eventService,
	}
	authorizationStore := &pq.AuthorizationStore{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	storeRedisLogger := idpsession.NewStoreRedisLogger(factory)
	idpsessionStoreRedis := &idpsession.StoreRedis{
		Redis:  appredisHandle,
		AppID:  appID,
		Clock:  clockClock,
		Logger: storeRedisLogger,
	}
	sessionConfig := appConfig.Session
	cookieManager := deps.NewCookieManager(request, trustProxy, httpConfig)
	cookieDef := session.NewSessionCookieDef(sessionConfig)
	idpsessionManager := &idpsession.Manager{
		Store:     idpsessionStoreRedis,
		Clock:     clockClock,
		Config:    sessionConfig,
		Cookies:   cookieManager,
		CookieDef: cookieDef,
	}
	redisLogger := redis.NewLogger(factory)
	redisStore := &redis.Store{
		Context:     contextContext,
		Redis:       appredisHandle,
		AppID:       appID,
		Logger:      redisLogger,
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
		Clock:       clockClock,
	}
	oAuthConfig := appConfig.OAuth
	sessionManager := &oauth2.SessionManager{
		Store:  redisStore,
		Clock:  clockClock,
		Config: oAuthConfig,
	}
	accountDeletionConfig := appConfig.AccountDeletion
	coordinator := &facade.Coordinator{
		Events:                eventService,
		Identities:            serviceService,
		Authenticators:        service4,
		Verification:          verificationService,
		MFA:                   mfaService,
		UserCommands:          commands,
		UserQueries:           queries,
		StdAttrsService:       stdattrsService,
		PasswordHistory:       historyStore,
		OAuth:                 authorizationStore,
		IDPSessions:           idpsessionManager,
		OAuthSessions:         sessionManager,
		IdentityConfig:        identityConfig,
		AccountDeletionConfig: accountDeletionConfig,
		Clock:                 clockClock,
	}
	userFacade := &facade.UserFacade{
		UserProvider: userProvider,
		Coordinator:  coordinator,
	}
	interactionLogger := interaction.NewLogger(factory)
	identityFacade := facade.IdentityFacade{
		Coordinator: coordinator,
	}
	authenticatorFacade := facade.AuthenticatorFacade{
		Coordinator: coordinator,
	}
	webEndpoints := &WebEndpoints{}
	hardSMSBucketer := &usage.HardSMSBucketer{
		FeatureConfig: featureConfig,
	}
	messageSender := &otp.MessageSender{
		Translation:     translationService,
		Endpoints:       webEndpoints,
		TaskQueue:       queue,
		Events:          eventService,
		RateLimiter:     limiter,
		HardSMSBucketer: hardSMSBucketer,
	}
	codeSender := &oob.CodeSender{
		OTPMessageSender: messageSender,
	}
	oAuthClientCredentials := deps.ProvideOAuthClientCredentials(secretConfig)
	normalizer := &stdattrs2.Normalizer{
		LoginIDNormalizerFactory: normalizerFactory,
	}
	oAuthProviderFactory := &sso.OAuthProviderFactory{
		Endpoints:                    webEndpoints,
		IdentityConfig:               identityConfig,
		Credentials:                  oAuthClientCredentials,
		RedirectURL:                  webEndpoints,
		Clock:                        clockClock,
		WechatURLProvider:            webEndpoints,
		StandardAttributesNormalizer: normalizer,
	}
	forgotPasswordConfig := appConfig.ForgotPassword
	forgotpasswordStore := &forgotpassword.Store{
		Context: contextContext,
		AppID:   appID,
		Redis:   appredisHandle,
	}
	providerLogger := forgotpassword.NewProviderLogger(factory)
	forgotpasswordProvider := &forgotpassword.Provider{
		RemoteIP:        remoteIP,
		Translation:     translationService,
		Config:          forgotPasswordConfig,
		Store:           forgotpasswordStore,
		Clock:           clockClock,
		URLs:            webEndpoints,
		TaskQueue:       queue,
		Logger:          providerLogger,
		Identities:      identityFacade,
		Authenticators:  authenticatorFacade,
		FeatureConfig:   featureConfig,
		Events:          eventService,
		RateLimiter:     limiter,
		HardSMSBucketer: hardSMSBucketer,
	}
	verificationCodeSender := &verification.CodeSender{
		OTPMessageSender: messageSender,
	}
	responseWriter := p.ResponseWriter
	nonceService := &nonce.Service{
		Cookies:        cookieManager,
		Request:        request,
		ResponseWriter: responseWriter,
	}
	challengeProvider := &challenge.Provider{
		Redis: appredisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	authenticationinfoStoreRedis := &authenticationinfo.StoreRedis{
		Context: contextContext,
		Redis:   appredisHandle,
		AppID:   appID,
	}
	eventStoreRedis := &access.EventStoreRedis{
		Redis: appredisHandle,
		AppID: appID,
	}
	eventProvider := &access.EventProvider{
		Store: eventStoreRedis,
	}
	rand := _wireRandValue
	idpsessionProvider := &idpsession.Provider{
		Context:         contextContext,
		RemoteIP:        remoteIP,
		UserAgentString: userAgentString,
		AppID:           appID,
		Redis:           appredisHandle,
		Store:           idpsessionStoreRedis,
		AccessEvents:    eventProvider,
		TrustProxy:      trustProxy,
		Config:          sessionConfig,
		Clock:           clockClock,
		Random:          rand,
	}
	mfaCookieDef := mfa.NewDeviceTokenCookieDef(authenticationConfig)
	whatsappStoreRedis := &whatsapp.StoreRedis{
		Context: contextContext,
		Redis:   appredisHandle,
		Clock:   clockClock,
	}
	whatsappLogger := whatsapp.NewLogger(factory)
	watiCredentials := deps.ProvideWATICredentials(secretConfig)
	whatsappProvider := &whatsapp.Provider{
		CodeStore:       whatsappStoreRedis,
		Clock:           clockClock,
		Logger:          whatsappLogger,
		WATICredentials: watiCredentials,
		Events:          eventService,
	}
	interactionContext := &interaction.Context{
		Request:                   request,
		RemoteIP:                  remoteIP,
		Database:                  sqlExecutor,
		Clock:                     clockClock,
		Config:                    appConfig,
		FeatureConfig:             featureConfig,
		Identities:                identityFacade,
		Authenticators:            authenticatorFacade,
		AnonymousIdentities:       anonymousProvider,
		BiometricIdentities:       biometricProvider,
		OOBAuthenticators:         oobProvider,
		OOBCodeSender:             codeSender,
		OAuthProviderFactory:      oAuthProviderFactory,
		MFA:                       mfaService,
		ForgotPassword:            forgotpasswordProvider,
		ResetPassword:             forgotpasswordProvider,
		LoginIDNormalizerFactory:  normalizerFactory,
		Verification:              verificationService,
		VerificationCodeSender:    verificationCodeSender,
		RateLimiter:               limiter,
		Nonces:                    nonceService,
		Challenges:                challengeProvider,
		Users:                     userProvider,
		StdAttrsService:           stdattrsService,
		Events:                    eventService,
		CookieManager:             cookieManager,
		AuthenticationInfoService: authenticationinfoStoreRedis,
		Sessions:                  idpsessionProvider,
		SessionManager:            idpsessionManager,
		SessionCookie:             cookieDef,
		MFADeviceTokenCookie:      mfaCookieDef,
		WhatsappCodeProvider:      whatsappProvider,
	}
	interactionStoreRedis := &interaction.StoreRedis{
		Redis: appredisHandle,
		AppID: appID,
	}
	interactionService := &interaction.Service{
		Logger:  interactionLogger,
		Context: interactionContext,
		Store:   interactionStoreRedis,
	}
	serviceInteractionService := &service3.InteractionService{
		Graph: interactionService,
	}
	facadeUserFacade := &facade2.UserFacade{
		UserSearchService:  elasticsearchService,
		Users:              userFacade,
		StandardAttributes: serviceNoEvent,
		Interaction:        serviceInteractionService,
	}
	auditLogFeatureConfig := featureConfig.AuditLog
	auditLogFacade := &facade2.AuditLogFacade{
		AuditLogQuery:         query,
		Clock:                 clockClock,
		AuditLogFeatureConfig: auditLogFeatureConfig,
	}
	facadeIdentityFacade := &facade2.IdentityFacade{
		Identities:  serviceService,
		Interaction: serviceInteractionService,
	}
	facadeAuthenticatorFacade := &facade2.AuthenticatorFacade{
		Authenticators: service4,
		Interaction:    serviceInteractionService,
	}
	verificationFacade := &facade2.VerificationFacade{
		Verification: verificationService,
	}
	manager2 := &session.Manager{
		IDPSessions:         idpsessionManager,
		AccessTokenSessions: sessionManager,
		Events:              eventService,
	}
	sessionFacade := &facade2.SessionFacade{
		Sessions: manager2,
	}
	userProfileFacade := &facade2.UserProfileFacade{
		StandardAttributes: serviceNoEvent,
		CustomAttributes:   customattrsServiceNoEvent,
		Events:             eventService,
	}
	graphqlContext := &graphql.Context{
		GQLLogger:           logger,
		Users:               userLoader,
		Identities:          identityLoader,
		Authenticators:      authenticatorLoader,
		AuditLogs:           auditLogLoader,
		UserFacade:          facadeUserFacade,
		AuditLogFacade:      auditLogFacade,
		IdentityFacade:      facadeIdentityFacade,
		AuthenticatorFacade: facadeAuthenticatorFacade,
		VerificationFacade:  verificationFacade,
		SessionFacade:       sessionFacade,
		UserProfileFacade:   userProfileFacade,
	}
	graphQLHandler := &transport.GraphQLHandler{
		GraphQLContext: graphqlContext,
		AppDatabase:    handle,
		AuditDatabase:  readHandle,
	}
	return graphQLHandler
}

var (
	_wireRandValue = idpsession.Rand(rand.SecureRand)
)

func newPresignImagesUploadHandler(p *deps.RequestProvider) http.Handler {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	jsonResponseWriterLogger := httputil.NewJSONResponseWriterLogger(factory)
	jsonResponseWriter := &httputil.JSONResponseWriter{
		Logger: jsonResponseWriterLogger,
	}
	request := p.Request
	rootProvider := appProvider.RootProvider
	environmentConfig := rootProvider.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	httpProto := deps.ProvideHTTPProto(request, trustProxy)
	httpHost := deps.ProvideHTTPHost(request, trustProxy)
	configConfig := appProvider.Config
	appConfig := configConfig.AppConfig
	appID := appConfig.ID
	secretConfig := configConfig.SecretConfig
	imagesKeyMaterials := deps.ProvideImagesKeyMaterials(secretConfig)
	clockClock := _wireSystemClockValue
	provider := &presign.Provider{
		Secret: imagesKeyMaterials,
		Clock:  clockClock,
		Host:   httpHost,
	}
	presignImagesUploadHandlerLogger := transport.NewPresignImagesUploadHandlerLogger(factory)
	presignImagesUploadHandler := &transport.PresignImagesUploadHandler{
		JSON:            jsonResponseWriter,
		HTTPProto:       httpProto,
		HTTPHost:        httpHost,
		AppID:           appID,
		PresignProvider: provider,
		Logger:          presignImagesUploadHandlerLogger,
	}
	return presignImagesUploadHandler
}
