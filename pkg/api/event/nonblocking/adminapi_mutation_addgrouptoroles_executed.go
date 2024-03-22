package nonblocking

import (
	"github.com/authgear/authgear-server/pkg/api/event"
)

const (
	AdminAPIMutationAddGroupToRolesExecuted event.Type = "admin_api.mutation.add_group_to_roles.executed"
)

type AdminAPIMutationAddGroupToRolesExecutedEventPayload struct {
	AffectedUserIDs []string `json:"-"`
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) NonBlockingEventType() event.Type {
	return AdminAPIMutationAddGroupToRolesExecuted
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) UserID() string {
	return ""
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) GetTriggeredBy() event.TriggeredByType {
	return event.TriggeredByTypeAdminAPI
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) FillContext(ctx *event.Context) {
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) ForHook() bool {
	return false
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) ForAudit() bool {
	// FIXME(tung): Should be true
	return false
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) RequireReindexUserIDs() []string {
	return e.AffectedUserIDs
}

func (e *AdminAPIMutationAddGroupToRolesExecutedEventPayload) DeletedUserIDs() []string {
	return []string{}
}

var _ event.NonBlockingPayload = &AdminAPIMutationAddGroupToRolesExecutedEventPayload{}
