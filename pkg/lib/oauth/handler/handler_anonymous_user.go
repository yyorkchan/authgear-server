package handler

import (
	"errors"
	"net/http"

	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/api/model"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticationinfo"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/interaction"
	interactionintents "github.com/authgear/authgear-server/pkg/lib/interaction/intents"
	"github.com/authgear/authgear-server/pkg/lib/interaction/nodes"
	"github.com/authgear/authgear-server/pkg/lib/oauth"
	"github.com/authgear/authgear-server/pkg/lib/oauth/protocol"
	"github.com/authgear/authgear-server/pkg/lib/session"
	"github.com/authgear/authgear-server/pkg/util/accesscontrol"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/log"
)

type anonymousSignupWithoutKeyInput struct{}

func (i *anonymousSignupWithoutKeyInput) GetAnonymousRequestToken() string { return "" }

func (i *anonymousSignupWithoutKeyInput) SignUpAnonymousUserWithoutKey() bool { return true }

var _ nodes.InputUseIdentityAnonymous = &anonymousSignupWithoutKeyInput{}

type AnonymousUserHandlerLogger struct{ *log.Logger }

func NewAnonymousUserHandlerLogger(lf *log.Factory) AnonymousUserHandlerLogger {
	return AnonymousUserHandlerLogger{lf.New("oauth-anonymous-user")}
}

type UserProvider interface {
	Get(id string, role accesscontrol.Role) (*model.User, error)
}

type CookiesGetter interface {
	GetCookies() []*http.Cookie
}

type SignupAnonymousUserResult struct {
	TokenResponse interface{}
	Cookies       []*http.Cookie
}

type AnonymousUserHandler struct {
	AppID       config.AppID
	OAuthConfig *config.OAuthConfig
	Logger      AnonymousUserHandlerLogger

	Graphs         GraphService
	Authorizations oauth.AuthorizationStore
	Clock          clock.Clock
	TokenService   TokenService
	UserProvider   UserProvider
}

// SignupAnonymousUser return token response or api errors
func (h *AnonymousUserHandler) SignupAnonymousUser(
	req *http.Request,
	clientID string,
	sessionType WebSessionType,
	refreshToken string,
) (*SignupAnonymousUserResult, error) {
	client, ok := h.OAuthConfig.GetClient(clientID)
	if !ok {
		// "invalid_client"
		return nil, apierrors.NewInvalid("invalid client ID")
	}

	if !*client.IsFirstParty {
		// unauthorized_client
		return nil, apierrors.NewInvalid("third-party clients may not use anonymous user")
	}

	switch sessionType {
	case WebSessionTypeCookie:
		return h.signupAnonymousUserWithCookieSessionType(req)
	case WebSessionTypeRefreshToken:
		return h.signupAnonymousUserWithRefreshTokenSessionType(client, refreshToken)
	default:
		panic("unknown web session type")
	}
}

func (h *AnonymousUserHandler) signupAnonymousUserWithCookieSessionType(
	req *http.Request,
) (*SignupAnonymousUserResult, error) {
	s := session.GetSession(req.Context())
	if s != nil && s.SessionType() == session.TypeIdentityProvider {
		user, err := h.UserProvider.Get(s.GetAuthenticationInfo().UserID, accesscontrol.EmptyRole)
		if err != nil {
			return nil, err
		}

		if user.IsAnonymous {
			return &SignupAnonymousUserResult{}, nil
		}
		return nil, apierrors.NewInvalid("user logged in as normal user, please logout first")
	}

	graph, err := h.runSignupAnonymousUserGraph(false)
	if err != nil {
		return nil, err
	}

	cookies := []*http.Cookie{}
	for _, node := range graph.Nodes {
		if a, ok := node.(CookiesGetter); ok {
			cookies = append(cookies, a.GetCookies()...)
		}
	}

	return &SignupAnonymousUserResult{
		Cookies: cookies,
	}, nil
}

func (h *AnonymousUserHandler) signupAnonymousUserWithRefreshTokenSessionType(
	client *config.OAuthClientConfig,
	refreshToken string,
) (*SignupAnonymousUserResult, error) {
	// TODO(oauth): allow specifying scopes for anonymous user signup
	scopes := []string{"openid", oauth.FullAccessScope}

	if refreshToken != "" {
		authz, grant, err := h.TokenService.ParseRefreshToken(refreshToken)
		if errors.Is(err, errInvalidRefreshToken) {
			return nil, apierrors.NewInvalid("invalid refresh token")
		} else if err != nil {
			return nil, err
		}

		resp := protocol.TokenResponse{}
		err = h.TokenService.IssueAccessGrant(client, scopes, authz.ID, authz.UserID,
			grant.ID, oauth.GrantSessionKindOffline, resp)
		if err != nil {
			return nil, err
		}

		return &SignupAnonymousUserResult{
			TokenResponse: resp,
		}, nil
	}

	graph, err := h.runSignupAnonymousUserGraph(true)
	if err != nil {
		return nil, err
	}

	info := authenticationinfo.T{
		UserID:          graph.MustGetUserID(),
		AuthenticatedAt: h.Clock.NowUTC(),
	}

	authz, err := checkAuthorization(
		h.Authorizations,
		h.Clock.NowUTC(),
		h.AppID,
		client.ClientID,
		info.UserID,
		scopes,
	)
	if err != nil {
		return nil, err
	}

	resp := protocol.TokenResponse{}
	opts := IssueOfflineGrantOptions{
		Scopes:             scopes,
		AuthorizationID:    authz.ID,
		AuthenticationInfo: info,
		DeviceInfo:         nil,
	}
	offlineGrant, err := h.TokenService.IssueOfflineGrant(client, opts, resp)
	if err != nil {
		return nil, err
	}

	err = h.TokenService.IssueAccessGrant(client, scopes, authz.ID, authz.UserID,
		offlineGrant.ID, oauth.GrantSessionKindOffline, resp)
	if err != nil {
		return nil, err
	}

	return &SignupAnonymousUserResult{
		TokenResponse: resp,
	}, nil
}

func (h *AnonymousUserHandler) runSignupAnonymousUserGraph(
	suppressIDPSessionCookie bool,
) (*interaction.Graph, error) {
	var graph *interaction.Graph
	err := h.Graphs.DryRun("", func(ctx *interaction.Context) (*interaction.Graph, error) {
		var err error
		intent := &interactionintents.IntentAuthenticate{
			Kind:                     interactionintents.IntentAuthenticateKindLogin,
			SuppressIDPSessionCookie: suppressIDPSessionCookie,
		}
		graph, err = h.Graphs.NewGraph(ctx, intent)
		if err != nil {
			return nil, err
		}

		var edges []interaction.Edge
		graph, edges, err = h.Graphs.Accept(ctx, graph, &anonymousSignupWithoutKeyInput{})
		if len(edges) != 0 {
			return nil, errors.New("interaction not completed for anonymous users")
		} else if err != nil {
			return nil, err
		}

		return graph, nil
	})

	if apierrors.IsKind(err, interaction.InvariantViolated) &&
		apierrors.AsAPIError(err).HasCause("AnonymousUserDisallowed") {
		// unauthorized_client
		return nil, apierrors.NewInvalid("anonymous user disallowed")
	} else if errors.Is(err, interaction.ErrInvalidCredentials) {
		// invalid_grant
		return nil, apierrors.NewInvalid(interaction.InvalidCredentials.Reason)
	} else if err != nil {
		return nil, err
	}

	err = h.Graphs.Run("", graph)
	if apierrors.IsAPIError(err) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return graph, nil
}
