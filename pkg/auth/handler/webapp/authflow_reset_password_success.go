package webapp

import (
	"net/http"

	"github.com/authgear/authgear-server/pkg/auth/handler/webapp/viewmodels"
	"github.com/authgear/authgear-server/pkg/auth/webapp"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/template"
)

var TemplateWebAuthflowResetPasswordSuccessHTML = template.RegisterHTML(
	"web/authflow_reset_password_success.html",
	components...,
)

func ConfigureAuthflowResetPasswordSuccessRoute(route httproute.Route) httproute.Route {
	return route.
		WithMethods("OPTIONS", "POST", "GET").
		WithPathPattern(webapp.AuthflowRouteResetPasswordSuccess)
}

type AuthflowResetPasswordSuccessHandler struct {
	Controller    *AuthflowController
	BaseViewModel *viewmodels.BaseViewModeler
	Renderer      Renderer
}

func (h *AuthflowResetPasswordSuccessHandler) GetData(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	baseViewModel := h.BaseViewModel.ViewModelForAuthFlow(r, w)
	viewmodels.Embed(data, baseViewModel)

	return data, nil
}

func (h *AuthflowResetPasswordSuccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handlers AuthflowControllerHandlers
	handlers.Get(func(s *webapp.Session, _ *webapp.AuthflowScreenWithFlowResponse) error {
		data, err := h.GetData(w, r)
		if err != nil {
			return err
		}

		h.Renderer.RenderHTML(w, r, TemplateWebAuthflowResetPasswordSuccessHTML, data)
		return nil
	})

	h.Controller.HandleWithoutFlow(w, r, &handlers)
}
