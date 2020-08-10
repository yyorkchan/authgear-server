package nodes

import (
	"errors"

	"github.com/authgear/authgear-server/pkg/auth/dependency/authenticator"
	"github.com/authgear/authgear-server/pkg/auth/dependency/newinteraction"
)

func init() {
	newinteraction.RegisterNode(&NodeCreateAuthenticatorTOTP{})
}

type InputCreateAuthenticatorTOTP interface {
	GetTOTP() string
	GetTOTPDisplayName() string
}

type EdgeCreateAuthenticatorTOTP struct {
	Stage         newinteraction.AuthenticationStage
	Authenticator *authenticator.Info
}

func (e *EdgeCreateAuthenticatorTOTP) Instantiate(ctx *newinteraction.Context, graph *newinteraction.Graph, rawInput interface{}) (newinteraction.Node, error) {
	input, ok := rawInput.(InputCreateAuthenticatorTOTP)
	if !ok {
		return nil, newinteraction.ErrIncompatibleInput
	}

	info := cloneAuthenticator(e.Authenticator)
	info.Props[authenticator.AuthenticatorPropTOTPDisplayName] = input.GetTOTPDisplayName()

	err := ctx.Authenticators.VerifySecret(info, nil, input.GetTOTP())
	if errors.Is(err, authenticator.ErrInvalidCredentials) {
		return nil, newinteraction.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	return &NodeCreateAuthenticatorTOTP{Stage: e.Stage, Authenticator: info}, nil
}

type NodeCreateAuthenticatorTOTP struct {
	Stage         newinteraction.AuthenticationStage `json:"stage"`
	Authenticator *authenticator.Info                `json:"authenticator"`
}

func (n *NodeCreateAuthenticatorTOTP) Prepare(ctx *newinteraction.Context, graph *newinteraction.Graph) error {
	return nil
}

func (n *NodeCreateAuthenticatorTOTP) Apply(perform func(eff newinteraction.Effect) error, graph *newinteraction.Graph) error {
	return nil
}

func (n *NodeCreateAuthenticatorTOTP) DeriveEdges(graph *newinteraction.Graph) ([]newinteraction.Edge, error) {
	return []newinteraction.Edge{
		&EdgeCreateAuthenticatorEnd{
			Stage:          n.Stage,
			Authenticators: []*authenticator.Info{n.Authenticator},
		},
	}, nil
}
