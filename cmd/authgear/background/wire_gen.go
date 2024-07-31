// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package background

import (
	"context"
	"github.com/authgear/authgear-server/pkg/lib/audit"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/oob"
	passkey3 "github.com/authgear/authgear-server/pkg/lib/authn/authenticator/passkey"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/password"
	service2 "github.com/authgear/authgear-server/pkg/lib/authn/authenticator/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/totp"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/anonymous"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/biometric"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/ldap"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/loginid"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/oauth"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/passkey"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/siwe"
	"github.com/authgear/authgear-server/pkg/lib/authn/mfa"
	"github.com/authgear/authgear-server/pkg/lib/authn/otp"
	"github.com/authgear/authgear-server/pkg/lib/authn/user"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/config/configsource"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/elasticsearch"
	"github.com/authgear/authgear-server/pkg/lib/endpoints"
	"github.com/authgear/authgear-server/pkg/lib/event"
	"github.com/authgear/authgear-server/pkg/lib/facade"
	"github.com/authgear/authgear-server/pkg/lib/feature/accountanonymization"
	"github.com/authgear/authgear-server/pkg/lib/feature/accountdeletion"
	"github.com/authgear/authgear-server/pkg/lib/feature/customattrs"
	passkey2 "github.com/authgear/authgear-server/pkg/lib/feature/passkey"
	siwe2 "github.com/authgear/authgear-server/pkg/lib/feature/siwe"
	"github.com/authgear/authgear-server/pkg/lib/feature/stdattrs"
	"github.com/authgear/authgear-server/pkg/lib/feature/verification"
	"github.com/authgear/authgear-server/pkg/lib/feature/web3"
	"github.com/authgear/authgear-server/pkg/lib/hook"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/appdb"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/auditdb"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/globaldb"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
	"github.com/authgear/authgear-server/pkg/lib/infra/redisqueue"
	"github.com/authgear/authgear-server/pkg/lib/lockout"
	oauth2 "github.com/authgear/authgear-server/pkg/lib/oauth"
	"github.com/authgear/authgear-server/pkg/lib/oauth/pq"
	"github.com/authgear/authgear-server/pkg/lib/oauth/redis"
	"github.com/authgear/authgear-server/pkg/lib/oauthclient"
	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/lib/rolesgroups"
	"github.com/authgear/authgear-server/pkg/lib/session"
	"github.com/authgear/authgear-server/pkg/lib/session/access"
	"github.com/authgear/authgear-server/pkg/lib/session/idpsession"
	"github.com/authgear/authgear-server/pkg/lib/translation"
	"github.com/authgear/authgear-server/pkg/lib/web"
	"github.com/authgear/authgear-server/pkg/util/backgroundjob"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/authgear/authgear-server/pkg/util/rand"
	"github.com/authgear/authgear-server/pkg/util/template"
)

// Injectors from wire.go:

func newConfigSourceController(p *deps.BackgroundProvider, c context.Context) *configsource.Controller {
	config := p.ConfigSourceConfig
	factory := p.LoggerFactory
	localFSLogger := configsource.NewLocalFSLogger(factory)
	manager := p.BaseResources
	localFS := &configsource.LocalFS{
		Logger:        localFSLogger,
		BaseResources: manager,
		Config:        config,
	}
	databaseLogger := configsource.NewDatabaseLogger(factory)
	environmentConfig := p.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	clock := _wireSystemClockValue
	globalDatabaseCredentialsEnvironmentConfig := &environmentConfig.GlobalDatabase
	sqlBuilder := globaldb.NewSQLBuilder(globalDatabaseCredentialsEnvironmentConfig)
	storeFactory := configsource.NewStoreFactory(c, sqlBuilder)
	pool := p.DatabasePool
	databaseEnvironmentConfig := &environmentConfig.DatabaseConfig
	databaseHandleFactory := configsource.NewDatabaseHandleFactory(c, pool, globalDatabaseCredentialsEnvironmentConfig, databaseEnvironmentConfig, factory)
	resolveAppIDType := configsource.NewResolveAppIDTypeDomain()
	database := &configsource.Database{
		Logger:                databaseLogger,
		BaseResources:         manager,
		TrustProxy:            trustProxy,
		Config:                config,
		Clock:                 clock,
		StoreFactory:          storeFactory,
		DatabaseHandleFactory: databaseHandleFactory,
		DatabaseCredentials:   globalDatabaseCredentialsEnvironmentConfig,
		DatabaseConfig:        databaseEnvironmentConfig,
		ResolveAppIDType:      resolveAppIDType,
	}
	controller := configsource.NewController(config, localFS, database)
	return controller
}

var (
	_wireSystemClockValue = clock.NewSystemClock()
)

func newAccountAnonymizationRunner(p *deps.BackgroundProvider, c context.Context, ctrl *configsource.Controller) *backgroundjob.Runner {
	factory := p.LoggerFactory
	pool := p.DatabasePool
	environmentConfig := p.EnvironmentConfig
	globalDatabaseCredentialsEnvironmentConfig := &environmentConfig.GlobalDatabase
	databaseEnvironmentConfig := &environmentConfig.DatabaseConfig
	clockClock := _wireSystemClockValue
	accountAnonymizationServiceFactory := &AccountAnonymizationServiceFactory{
		BackgroundProvider: p,
	}
	runnableFactory := accountanonymization.NewRunnableFactory(pool, globalDatabaseCredentialsEnvironmentConfig, databaseEnvironmentConfig, factory, clockClock, ctrl, accountAnonymizationServiceFactory)
	runner := accountanonymization.NewRunner(factory, runnableFactory)
	return runner
}

func newAccountDeletionRunner(p *deps.BackgroundProvider, c context.Context, ctrl *configsource.Controller) *backgroundjob.Runner {
	factory := p.LoggerFactory
	pool := p.DatabasePool
	environmentConfig := p.EnvironmentConfig
	globalDatabaseCredentialsEnvironmentConfig := &environmentConfig.GlobalDatabase
	databaseEnvironmentConfig := &environmentConfig.DatabaseConfig
	clockClock := _wireSystemClockValue
	accountDeletionServiceFactory := &AccountDeletionServiceFactory{
		BackgroundProvider: p,
	}
	runnableFactory := accountdeletion.NewRunnableFactory(pool, globalDatabaseCredentialsEnvironmentConfig, databaseEnvironmentConfig, factory, clockClock, ctrl, accountDeletionServiceFactory)
	runner := accountdeletion.NewRunner(factory, runnableFactory)
	return runner
}

func newUserService(ctx context.Context, p *deps.BackgroundProvider, appID string, appContext *config.AppContext) *UserService {
	pool := p.DatabasePool
	environmentConfig := p.EnvironmentConfig
	databaseEnvironmentConfig := &environmentConfig.DatabaseConfig
	configConfig := appContext.Config
	secretConfig := configConfig.SecretConfig
	databaseCredentials := deps.ProvideDatabaseCredentials(secretConfig)
	factory := p.LoggerFactory
	handle := appdb.NewHandle(ctx, pool, databaseEnvironmentConfig, databaseCredentials, factory)
	appConfig := configConfig.AppConfig
	configAppID := appConfig.ID
	sqlBuilderApp := appdb.NewSQLBuilderApp(databaseCredentials, configAppID)
	sqlExecutor := appdb.NewSQLExecutor(ctx, handle)
	clockClock := _wireSystemClockValue
	store := &user.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
		Clock:       clockClock,
		AppID:       configAppID,
	}
	rawCommands := &user.RawCommands{
		Store: store,
		Clock: clockClock,
	}
	rawQueries := &user.RawQueries{
		Store: store,
	}
	remoteIP := ProvideRemoteIP()
	userAgentString := ProvideUserAgentString()
	logger := event.NewLogger(factory)
	localizationConfig := appConfig.Localization
	sqlBuilder := appdb.NewSQLBuilder(databaseCredentials)
	storeImpl := event.NewStoreImpl(sqlBuilder, sqlExecutor)
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
	uiConfig := appConfig.UI
	manager := p.BaseResources
	typeCheckerFactory := &loginid.TypeCheckerFactory{
		UIConfig:      uiConfig,
		LoginIDConfig: loginIDConfig,
		Resources:     manager,
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
	redisPool := p.RedisPool
	hub := p.RedisHub
	redisEnvironmentConfig := &environmentConfig.RedisConfig
	redisCredentials := deps.ProvideRedisCredentials(secretConfig)
	appredisHandle := appredis.NewHandle(redisPool, hub, redisEnvironmentConfig, redisCredentials, factory)
	store2 := &passkey2.Store{
		Context: ctx,
		Redis:   appredisHandle,
		AppID:   configAppID,
	}
	request := NewDummyHTTPRequest()
	trustProxy := environmentConfig.TrustProxy
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
	httpProto := ProvideHTTPProto()
	httpHost := ProvideHTTPHost()
	httpOrigin := httputil.MakeHTTPOrigin(httpProto, httpHost)
	webAppCDNHost := environmentConfig.WebAppCDNHost
	globalEmbeddedResourceManager := p.EmbeddedResources
	staticAssetResolver := &web.StaticAssetResolver{
		Context:           ctx,
		Localization:      localizationConfig,
		HTTPOrigin:        httpOrigin,
		HTTPProto:         httpProto,
		WebAppCDNHost:     webAppCDNHost,
		Resources:         manager,
		EmbeddedResources: globalEmbeddedResourceManager,
	}
	translationService := &translation.Service{
		Context:        ctx,
		TemplateEngine: engine,
		StaticAssets:   staticAssetResolver,
	}
	configService := &passkey2.ConfigService{
		Request:            request,
		TrustProxy:         trustProxy,
		TranslationService: translationService,
	}
	passkeyService := &passkey2.Service{
		Store:         store2,
		ConfigService: configService,
	}
	passkeyProvider := &passkey.Provider{
		Store:   passkeyStore,
		Clock:   clockClock,
		Passkey: passkeyService,
	}
	siweStore := &siwe.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	web3Config := appConfig.Web3
	storeRedis := &siwe2.StoreRedis{
		Context: ctx,
		Redis:   appredisHandle,
		AppID:   configAppID,
		Clock:   clockClock,
	}
	ratelimitLogger := ratelimit.NewLogger(factory)
	storageRedis := &ratelimit.StorageRedis{
		AppID: configAppID,
		Redis: appredisHandle,
	}
	rateLimitsFeatureConfig := featureConfig.RateLimits
	limiter := &ratelimit.Limiter{
		Logger:  ratelimitLogger,
		Storage: storageRedis,
		Config:  rateLimitsFeatureConfig,
	}
	siweLogger := siwe2.NewLogger(factory)
	siweService := &siwe2.Service{
		RemoteIP:             remoteIP,
		HTTPOrigin:           httpOrigin,
		Web3Config:           web3Config,
		AuthenticationConfig: authenticationConfig,
		Clock:                clockClock,
		NonceStore:           storeRedis,
		RateLimiter:          limiter,
		Logger:               siweLogger,
	}
	siweProvider := &siwe.Provider{
		Store: siweStore,
		Clock: clockClock,
		SIWE:  siweService,
	}
	ldapStore := &ldap.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	ldapProvider := &ldap.Provider{
		Store: ldapStore,
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
		SIWE:                  siweProvider,
		LDAP:                  ldapProvider,
	}
	store3 := &service2.Store{
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
	expiry := password.ProvideExpiry(authenticatorPasswordConfig, clockClock)
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
		Expiry:          expiry,
		Housekeeper:     housekeeper,
	}
	store4 := &passkey3.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	provider2 := &passkey3.Provider{
		Store:   store4,
		Clock:   clockClock,
		Passkey: passkeyService,
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
	oobStore := &oob.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	oobProvider := &oob.Provider{
		Store:                    oobStore,
		LoginIDNormalizerFactory: normalizerFactory,
		Clock:                    clockClock,
	}
	testModeConfig := appConfig.TestMode
	testModeFeatureConfig := featureConfig.TestMode
	codeStoreRedis := &otp.CodeStoreRedis{
		Redis: appredisHandle,
		AppID: configAppID,
		Clock: clockClock,
	}
	lookupStoreRedis := &otp.LookupStoreRedis{
		Redis: appredisHandle,
		AppID: configAppID,
		Clock: clockClock,
	}
	attemptTrackerRedis := &otp.AttemptTrackerRedis{
		Redis: appredisHandle,
		AppID: configAppID,
		Clock: clockClock,
	}
	otpLogger := otp.NewLogger(factory)
	otpService := &otp.Service{
		Clock:                 clockClock,
		AppID:                 configAppID,
		TestModeConfig:        testModeConfig,
		TestModeFeatureConfig: testModeFeatureConfig,
		RemoteIP:              remoteIP,
		CodeStore:             codeStoreRedis,
		LookupStore:           lookupStoreRedis,
		AttemptTracker:        attemptTrackerRedis,
		Logger:                otpLogger,
		RateLimiter:           limiter,
	}
	rateLimits := service2.RateLimits{
		IP:          remoteIP,
		Config:      authenticationConfig,
		RateLimiter: limiter,
	}
	authenticationLockoutConfig := authenticationConfig.Lockout
	lockoutLogger := lockout.NewLogger(factory)
	lockoutStorageRedis := &lockout.StorageRedis{
		AppID: configAppID,
		Redis: appredisHandle,
	}
	lockoutService := &lockout.Service{
		Logger:  lockoutLogger,
		Storage: lockoutStorageRedis,
	}
	serviceLockout := service2.Lockout{
		Config:   authenticationLockoutConfig,
		RemoteIP: remoteIP,
		Provider: lockoutService,
	}
	service3 := &service2.Service{
		Store:          store3,
		Config:         appConfig,
		Password:       passwordProvider,
		Passkey:        provider2,
		TOTP:           totpProvider,
		OOBOTP:         oobProvider,
		OTPCodeService: otpService,
		RateLimits:     rateLimits,
		Lockout:        serviceLockout,
	}
	verificationConfig := appConfig.Verification
	userProfileConfig := appConfig.UserProfile
	storePQ := &verification.StorePQ{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	verificationService := &verification.Service{
		Config:            verificationConfig,
		UserProfileConfig: userProfileConfig,
		Clock:             clockClock,
		ClaimStore:        storePQ,
	}
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
	nftIndexerAPIEndpoint := environmentConfig.NFTIndexerAPIEndpoint
	web3Service := &web3.Service{
		APIEndpoint: nftIndexerAPIEndpoint,
		Web3Config:  web3Config,
	}
	rolesgroupsStore := &rolesgroups.Store{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
		Clock:       clockClock,
	}
	queries := &rolesgroups.Queries{
		Store: rolesgroupsStore,
	}
	userQueries := &user.Queries{
		RawQueries:         rawQueries,
		Store:              store,
		Identities:         serviceService,
		Authenticators:     service3,
		Verification:       verificationService,
		StandardAttributes: serviceNoEvent,
		CustomAttributes:   customattrsServiceNoEvent,
		Web3:               web3Service,
		RolesAndGroups:     queries,
	}
	resolverImpl := &event.ResolverImpl{
		Users: userQueries,
	}
	hookLogger := hook.NewLogger(factory)
	hookConfig := appConfig.Hook
	webHookLogger := hook.NewWebHookLogger(factory)
	webhookKeyMaterials := deps.ProvideWebhookKeyMaterials(secretConfig)
	webHookImpl := hook.WebHookImpl{
		Logger: webHookLogger,
		Secret: webhookKeyMaterials,
	}
	syncHTTPClient := hook.NewSyncHTTPClient(hookConfig)
	asyncHTTPClient := hook.NewAsyncHTTPClient()
	eventWebHookImpl := &hook.EventWebHookImpl{
		WebHookImpl: webHookImpl,
		SyncHTTP:    syncHTTPClient,
		AsyncHTTP:   asyncHTTPClient,
	}
	denoHookLogger := hook.NewDenoHookLogger(factory)
	denoHook := hook.DenoHook{
		Context:         ctx,
		ResourceManager: manager,
		Logger:          denoHookLogger,
	}
	denoEndpoint := environmentConfig.DenoEndpoint
	syncDenoClient := hook.NewSyncDenoClient(denoEndpoint, hookConfig, hookLogger)
	asyncDenoClient := hook.NewAsyncDenoClient(denoEndpoint, hookLogger)
	eventDenoHookImpl := &hook.EventDenoHookImpl{
		DenoHook:        denoHook,
		SyncDenoClient:  syncDenoClient,
		AsyncDenoClient: asyncDenoClient,
	}
	commands := &rolesgroups.Commands{
		Store: rolesgroupsStore,
	}
	sink := &hook.Sink{
		Logger:             hookLogger,
		Config:             hookConfig,
		Clock:              clockClock,
		EventWebHook:       eventWebHookImpl,
		EventDenoHook:      eventDenoHookImpl,
		StandardAttributes: serviceNoEvent,
		CustomAttributes:   customattrsServiceNoEvent,
		RolesAndGroups:     commands,
	}
	auditLogger := audit.NewLogger(factory)
	auditDatabaseCredentials := deps.ProvideAuditDatabaseCredentials(secretConfig)
	writeHandle := auditdb.NewWriteHandle(ctx, pool, databaseEnvironmentConfig, auditDatabaseCredentials, factory)
	auditdbSQLBuilderApp := auditdb.NewSQLBuilderApp(auditDatabaseCredentials, configAppID)
	writeSQLExecutor := auditdb.NewWriteSQLExecutor(ctx, writeHandle)
	writeStore := &audit.WriteStore{
		SQLBuilder:  auditdbSQLBuilderApp,
		SQLExecutor: writeSQLExecutor,
	}
	auditSink := &audit.Sink{
		Logger:   auditLogger,
		Database: writeHandle,
		Store:    writeStore,
	}
	elasticsearchLogger := elasticsearch.NewLogger(factory)
	elasticsearchServiceLogger := elasticsearch.NewElasticsearchServiceLogger(factory)
	elasticsearchCredentials := deps.ProvideElasticsearchCredentials(secretConfig)
	client := elasticsearch.NewClient(elasticsearchCredentials)
	noopTaskQueue := NewNoopTaskQueue()
	userReindexProducer := redisqueue.NewUserReindexProducer(appredisHandle, clockClock)
	elasticsearchService := elasticsearch.Service{
		Clock:       clockClock,
		Context:     ctx,
		Database:    handle,
		Logger:      elasticsearchServiceLogger,
		AppID:       configAppID,
		Client:      client,
		Users:       userQueries,
		UserStore:   store,
		OAuth:       oauthStore,
		LoginID:     loginidStore,
		RolesGroups: rolesgroupsStore,
		TaskQueue:   noopTaskQueue,
		Producer:    userReindexProducer,
	}
	elasticsearchSink := &elasticsearch.Sink{
		Logger:   elasticsearchLogger,
		Service:  elasticsearchService,
		Database: handle,
	}
	eventService := event.NewService(ctx, configAppID, remoteIP, userAgentString, logger, handle, clockClock, localizationConfig, storeImpl, resolverImpl, sink, auditSink, elasticsearchSink)
	userCommands := &user.Commands{
		RawCommands:        rawCommands,
		RawQueries:         rawQueries,
		Events:             eventService,
		Verification:       verificationService,
		UserProfileConfig:  userProfileConfig,
		StandardAttributes: serviceNoEvent,
		CustomAttributes:   customattrsServiceNoEvent,
		Web3:               web3Service,
		RolesAndGroups:     queries,
	}
	userProvider := &user.Provider{
		Commands: userCommands,
		Queries:  userQueries,
	}
	storeDeviceTokenRedis := &mfa.StoreDeviceTokenRedis{
		Redis: appredisHandle,
		AppID: configAppID,
		Clock: clockClock,
	}
	storeRecoveryCodePQ := &mfa.StoreRecoveryCodePQ{
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
	}
	mfaLockout := mfa.Lockout{
		Config:   authenticationLockoutConfig,
		RemoteIP: remoteIP,
		Provider: lockoutService,
	}
	mfaService := &mfa.Service{
		IP:            remoteIP,
		DeviceTokens:  storeDeviceTokenRedis,
		RecoveryCodes: storeRecoveryCodePQ,
		Clock:         clockClock,
		Config:        authenticationConfig,
		RateLimiter:   limiter,
		Lockout:       mfaLockout,
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
		AppID:  configAppID,
		Clock:  clockClock,
		Logger: storeRedisLogger,
	}
	sessionConfig := appConfig.Session
	httpConfig := appConfig.HTTP
	cookieManager := deps.NewCookieManager(request, trustProxy, httpConfig)
	cookieDef := session.NewSessionCookieDef(sessionConfig)
	idpsessionManager := &idpsession.Manager{
		Store:     idpsessionStoreRedis,
		Config:    sessionConfig,
		Cookies:   cookieManager,
		CookieDef: cookieDef,
	}
	redisLogger := redis.NewLogger(factory)
	redisStore := &redis.Store{
		Context:     ctx,
		Redis:       appredisHandle,
		AppID:       configAppID,
		Logger:      redisLogger,
		SQLBuilder:  sqlBuilderApp,
		SQLExecutor: sqlExecutor,
		Clock:       clockClock,
	}
	oAuthConfig := appConfig.OAuth
	eventStoreRedis := &access.EventStoreRedis{
		Redis: appredisHandle,
		AppID: configAppID,
	}
	eventProvider := &access.EventProvider{
		Store: eventStoreRedis,
	}
	rand := _wireRandValue
	idpsessionProvider := &idpsession.Provider{
		Context:         ctx,
		RemoteIP:        remoteIP,
		UserAgentString: userAgentString,
		AppID:           configAppID,
		Redis:           appredisHandle,
		Store:           idpsessionStoreRedis,
		AccessEvents:    eventProvider,
		TrustProxy:      trustProxy,
		Config:          sessionConfig,
		Clock:           clockClock,
		Random:          rand,
	}
	endpointsEndpoints := &endpoints.Endpoints{
		HTTPHost:  httpHost,
		HTTPProto: httpProto,
	}
	oauthclientResolver := &oauthclient.Resolver{
		OAuthConfig:     oAuthConfig,
		TesterEndpoints: endpointsEndpoints,
	}
	offlineGrantService := oauth2.OfflineGrantService{
		OAuthConfig:    oAuthConfig,
		Clock:          clockClock,
		IDPSessions:    idpsessionProvider,
		ClientResolver: oauthclientResolver,
		OfflineGrants:  redisStore,
	}
	sessionManager := &oauth2.SessionManager{
		Store:   redisStore,
		Config:  oAuthConfig,
		Service: offlineGrantService,
	}
	accountDeletionConfig := appConfig.AccountDeletion
	accountAnonymizationConfig := appConfig.AccountAnonymization
	coordinator := &facade.Coordinator{
		Events:                     eventService,
		Identities:                 serviceService,
		Authenticators:             service3,
		Verification:               verificationService,
		MFA:                        mfaService,
		UserCommands:               userCommands,
		UserQueries:                userQueries,
		RolesGroupsCommands:        commands,
		StdAttrsService:            stdattrsService,
		PasswordHistory:            historyStore,
		OAuth:                      authorizationStore,
		IDPSessions:                idpsessionManager,
		OAuthSessions:              sessionManager,
		IdentityConfig:             identityConfig,
		AccountDeletionConfig:      accountDeletionConfig,
		AccountAnonymizationConfig: accountAnonymizationConfig,
		Clock:                      clockClock,
	}
	userFacade := &facade.UserFacade{
		UserProvider: userProvider,
		Coordinator:  coordinator,
	}
	userService := &UserService{
		AppDBHandle: handle,
		UserFacade:  userFacade,
	}
	return userService
}

var (
	_wireRandValue = idpsession.Rand(rand.SecureRand)
)
