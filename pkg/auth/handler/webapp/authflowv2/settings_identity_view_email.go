package authflowv2

import (
	"net/http"

	handlerwebapp "github.com/authgear/authgear-server/pkg/auth/handler/webapp"
	"github.com/authgear/authgear-server/pkg/auth/handler/webapp/viewmodels"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	identityservice "github.com/authgear/authgear-server/pkg/lib/authn/identity/service"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/feature/verification"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/appdb"
	"github.com/authgear/authgear-server/pkg/lib/session"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/template"
)

var TemplateWebSettingsIdentityViewEmailHTML = template.RegisterHTML(
	"web/authflowv2/settings_identity_view_email.html",
	handlerwebapp.SettingsComponents...,
)

func ConfigureAuthflowV2SettingsIdentityViewEmailRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("OPTIONS", "POST", "GET").
		WithPathPattern(AuthflowV2RouteSettingsIdentityViewEmail)
}

type AuthflowV2SettingsIdentityViewEmailViewModel struct {
	LoginIDKey     string
	EmailIdentity  *identity.LoginID
	Verifications  map[string][]verification.ClaimStatus
	UpdateDisabled bool
	DeleteDisabled bool
}

type AuthflowV2SettingsIdentityViewEmailHandler struct {
	Database          *appdb.Handle
	LoginIDConfig     *config.LoginIDConfig
	Identities        *identityservice.Service
	ControllerFactory handlerwebapp.ControllerFactory
	BaseViewModel     *viewmodels.BaseViewModeler
	Verification      handlerwebapp.SettingsVerificationService
	Renderer          handlerwebapp.Renderer
}

func (h *AuthflowV2SettingsIdentityViewEmailHandler) GetData(r *http.Request, rw http.ResponseWriter) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	loginIDKey := r.Form.Get("q_login_id_key")
	loginID := r.Form.Get("q_login_id")

	baseViewModel := h.BaseViewModel.ViewModel(r, rw)
	viewmodels.Embed(data, baseViewModel)

	userID := session.GetUserID(r.Context())

	emailIdentity, err := h.Identities.LoginID.Get(*userID, loginID)
	if err != nil {
		return nil, err
	}

	verifications, err := h.Verification.GetVerificationStatuses([]*identity.Info{emailIdentity.ToInfo()})
	if err != nil {
		return nil, err
	}

	updateDisabled := true
	deleteDisabled := true
	if loginIDConfig, ok := h.LoginIDConfig.GetKeyConfig(loginIDKey); ok {
		updateDisabled = *loginIDConfig.UpdateDisabled
		deleteDisabled = *loginIDConfig.DeleteDisabled
	}

	vm := AuthflowV2SettingsIdentityViewEmailViewModel{
		LoginIDKey:     loginIDKey,
		EmailIdentity:  emailIdentity,
		Verifications:  verifications,
		UpdateDisabled: updateDisabled,
		DeleteDisabled: deleteDisabled,
	}
	viewmodels.Embed(data, vm)

	return data, nil
}

func (h *AuthflowV2SettingsIdentityViewEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctrl, err := h.ControllerFactory.New(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer ctrl.ServeWithoutDBTx()

	ctrl.Get(func() error {
		var data map[string]interface{}
		err := h.Database.WithTx(func() error {
			data, err = h.GetData(r, w)
			return err
		})
		if err != nil {
			return err
		}

		h.Renderer.RenderHTML(w, r, TemplateWebSettingsIdentityViewEmailHTML, data)
		return nil
	})
}
