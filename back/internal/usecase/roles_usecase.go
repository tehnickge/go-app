package usecase

import (
	"errors"
	"strings"

	"backend-app/internal/domain"
	"backend-app/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRoleRequest struct {
	Name        string
	Description string
}

type UpdateRoleRequest struct {
	Name        *string
	Description *string
}

type RolesUsecase struct {
	rolesRepo *repository.RoleRepository
	authz     *Authorizer
}

func NewRolesUsecase(rolesRepo *repository.RoleRepository, authz *Authorizer) *RolesUsecase {
	return &RolesUsecase{rolesRepo: rolesRepo, authz: authz}
}

func (uc *RolesUsecase) Create(callerID uuid.UUID, req CreateRoleRequest) (*RoleDTO, error) {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("invalid input")
	}

	role := &domain.Role{
		Name:        name,
		Description: req.Description,
	}

	if err := uc.rolesRepo.Create(role); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}

	return toRoleDTO(role), nil
}

func (uc *RolesUsecase) GetByID(_ uuid.UUID, roleID uuid.UUID) (*RoleDTO, error) {
	role, err := uc.rolesRepo.GetByID(roleID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return toRoleDTO(role), nil
}

func (uc *RolesUsecase) List(_ uuid.UUID) ([]*RoleDTO, error) {
	roles, err := uc.rolesRepo.List()
	if err != nil {
		return nil, err
	}
	out := make([]*RoleDTO, 0, len(roles))
	for _, r := range roles {
		out = append(out, toRoleDTO(r))
	}
	return out, nil
}

func (uc *RolesUsecase) Update(callerID uuid.UUID, roleID uuid.UUID, req UpdateRoleRequest) (*RoleDTO, error) {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	fields := map[string]interface{}{}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("invalid input")
		}
		fields["name"] = name
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if len(fields) == 0 {
		return nil, errors.New("nothing to update")
	}

	if err := uc.rolesRepo.UpdateFields(roleID.String(), fields); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}

	role, err := uc.rolesRepo.GetByID(roleID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return toRoleDTO(role), nil
}

func (uc *RolesUsecase) SoftDelete(callerID uuid.UUID, roleID uuid.UUID) error {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrForbidden
	}
	if err := uc.rolesRepo.SoftDelete(roleID.String()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func toRoleDTO(r *domain.Role) *RoleDTO {
	return &RoleDTO{
		ID:          r.ID.String(),
		Name:        r.Name,
		Description: r.Description,
	}
}

