package blocking

import (
	"github.com/authgear/authgear-server/pkg/api/event"
	"github.com/authgear/authgear-server/pkg/api/model"
	"github.com/authgear/authgear-server/pkg/lib/rolesgroups"
	"github.com/authgear/authgear-server/pkg/util/accesscontrol"
)

func ApplyUserMutations(user model.User, userMutations event.UserMutations) (out model.User, mutated bool) {
	if userMutations.StandardAttributes != nil {
		user.StandardAttributes = userMutations.StandardAttributes
		mutated = true
	}
	if userMutations.CustomAttributes != nil {
		user.CustomAttributes = userMutations.CustomAttributes
		mutated = true
	}
	if userMutations.Roles != nil {
		user.Roles = userMutations.Roles
		mutated = true
	}
	if userMutations.Groups != nil {
		user.Groups = userMutations.Groups
		mutated = true
	}

	out = user
	return
}

func MakeUserMutations(user model.User) event.UserMutations {
	return event.UserMutations{
		StandardAttributes: user.StandardAttributes,
		CustomAttributes:   user.CustomAttributes,
		Roles:              user.Roles,
		Groups:             user.Groups,
	}
}

func PerformEffectsOnUser(ctx event.MutationsEffectContext, userID string, userMutations event.UserMutations) error {
	if userMutations.StandardAttributes != nil {
		err := ctx.StandardAttributes.UpdateStandardAttributes(
			accesscontrol.RoleGreatest,
			userID,
			userMutations.StandardAttributes,
		)
		if err != nil {
			return err
		}
	}
	if userMutations.CustomAttributes != nil {
		err := ctx.CustomAttributes.UpdateAllCustomAttributes(
			accesscontrol.RoleGreatest,
			userID,
			userMutations.CustomAttributes,
		)
		if err != nil {
			return err
		}
	}
	if userMutations.Roles != nil {
		err := ctx.RolesAndGroups.ResetUserRole(
			&rolesgroups.ResetUserRoleOptions{
				UserID:   userID,
				RoleKeys: userMutations.Roles,
			},
		)
		if err != nil {
			return err
		}
	}
	if userMutations.Groups != nil {
		err := ctx.RolesAndGroups.ResetUserGroup(
			&rolesgroups.ResetUserGroupOptions{
				UserID:    userID,
				GroupKeys: userMutations.Groups,
			},
		)
		if err != nil {
			return err
		}

	}

	return nil
}
