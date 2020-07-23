package nodes

import (
	"github.com/authgear/authgear-server/pkg/auth/dependency/authenticator"
	"github.com/authgear/authgear-server/pkg/auth/dependency/identity"
	"github.com/authgear/authgear-server/pkg/auth/dependency/newinteraction"
	"github.com/authgear/authgear-server/pkg/core/authn"
)

func getAuthenticators(
	ctx *newinteraction.Context,
	graph *newinteraction.Graph,
	stage newinteraction.AuthenticationStage,
	typ authn.AuthenticatorType,
) (*identity.Info, []*authenticator.Info, error) {
	var identityInfo *identity.Info
	var infos []*authenticator.Info
	var err error
	if stage == newinteraction.AuthenticationStagePrimary {
		identityInfo = graph.MustGetUserIdentity()
		infos, err = ctx.Authenticators.ListByIdentity(identityInfo.UserID, identityInfo)

		n := 0
		for _, info := range infos {
			if info.Type == typ {
				infos[n] = info
				n++
			}
		}
		infos = infos[:n]
	} else {
		userID := graph.MustGetUserID()
		infos, err = ctx.Authenticators.List(userID, typ)
	}
	if err != nil {
		return nil, nil, err
	}

	return identityInfo, infos, nil
}
