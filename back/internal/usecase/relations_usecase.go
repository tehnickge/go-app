package usecase

import (
	"backend-app/internal/repository"

	"github.com/google/uuid"
)

type RelationsUsecase struct {
	usersRolesRepo  *repository.UsersRolesRepository
	usersGroupsRepo *repository.UsersGroupsRepository
	authz            *Authorizer
}

func NewRelationsUsecase(
	usersRolesRepo *repository.UsersRolesRepository,
	usersGroupsRepo *repository.UsersGroupsRepository,
	authz *Authorizer,
) *RelationsUsecase {
	return &RelationsUsecase{
		usersRolesRepo:  usersRolesRepo,
		usersGroupsRepo: usersGroupsRepo,
		authz:            authz,
	}
}

func (uc *RelationsUsecase) AssignRoleToUser(callerID, userID, roleID uuid.UUID) error {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrForbidden
	}
	return uc.usersRolesRepo.AssignRoleToUser(userID, roleID)
}

func (uc *RelationsUsecase) RemoveRoleFromUser(callerID, userID, roleID uuid.UUID) error {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrForbidden
	}
	return uc.usersRolesRepo.RemoveRoleFromUser(userID, roleID)
}

type AddUserToGroupRequest struct {
	Role string `json:"role"`
}

func (uc *RelationsUsecase) AddUserToGroup(callerID, groupID, userID uuid.UUID, req AddUserToGroupRequest) error {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrForbidden
	}
	if req.Role == "" {
		req.Role = "member"
	}
	return uc.usersGroupsRepo.AddUserToGroup(userID, groupID, req.Role)
}

func (uc *RelationsUsecase) RemoveUserFromGroup(callerID, groupID, userID uuid.UUID) error {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrForbidden
	}
	return uc.usersGroupsRepo.RemoveUserFromGroup(userID, groupID)
}

func (uc *RelationsUsecase) GetUserRoles(callerID, userID uuid.UUID) ([]*RoleDTO, error) {
	if callerID != userID {
		isAdmin, err := uc.authz.IsAdmin(callerID)
		if err != nil {
			return nil, err
		}
		if !isAdmin {
			return nil, ErrForbidden
		}
	}
	roles, err := uc.usersRolesRepo.GetActiveRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	out := make([]*RoleDTO, 0, len(roles))
	for _, r := range roles {
		if r == nil {
			continue
		}
		out = append(out, &RoleDTO{ID: r.ID.String(), Name: r.Name, Description: r.Description})
	}
	return out, nil
}

func (uc *RelationsUsecase) GetUserGroups(callerID, userID uuid.UUID) ([]*GroupDTO, error) {
	if callerID != userID {
		isAdmin, err := uc.authz.IsAdmin(callerID)
		if err != nil {
			return nil, err
		}
		if !isAdmin {
			return nil, ErrForbidden
		}
	}

	groups, err := uc.usersGroupsRepo.GetActiveGroupsByUserID(userID)
	if err != nil {
		return nil, err
	}

	out := make([]*GroupDTO, 0, len(groups))
	for _, g := range groups {
		if g == nil {
			continue
		}
		dto := toGroupDTO(g)
		if dto != nil {
			out = append(out, dto)
		}
	}
	return out, nil
}

