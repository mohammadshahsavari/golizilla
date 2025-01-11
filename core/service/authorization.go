package service

import (
	"context"

	"github.com/google/uuid"
)

type IAuthorizationService interface {
	IsAuthorized(ctx context.Context, cuserCtx context.Context, userId uuid.UUID, requirePrivileges ...string) (bool, error)
	IsAuthorizedOnInsance(ctx context.Context, cuserCtx context.Context, userId uuid.UUID, questionnaireId uuid.UUID, requirePrivileges ...string) (bool, error)
}

type authorizationService struct {
	roleService IRoleService
}

func NewAuthorizationService(roleService IRoleService) IAuthorizationService {
	return &authorizationService{
		roleService: roleService,
	}
}

func (s *authorizationService) IsAuthorized(ctx context.Context, cuserCtx context.Context, userId uuid.UUID, requirePrivileges ...string) (bool, error) {
	return s.roleService.HasPrivileges(ctx, cuserCtx, userId, requirePrivileges...)
}

func (s *authorizationService) IsAuthorizedOnInsance(ctx context.Context, cuserCtx context.Context, userId uuid.UUID, questionnaireId uuid.UUID, requirePrivileges ...string) (bool, error) {
	return s.roleService.HasPrivilegesOnInsance(ctx, cuserCtx, userId, questionnaireId, requirePrivileges...)
}
