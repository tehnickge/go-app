package usecase

import (
	"errors"
	"strings"

	"backend-app/internal/domain"
	"backend-app/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	OwnerID     *string `json:"ownerId,omitempty"`
}

type CreateGroupRequest struct {
	Name        string
	Description string
}

type UpdateGroupRequest struct {
	Name        *string
	Description *string
}

type GroupsUsecase struct {
	groupsRepo *repository.GroupRepository
	authz       *Authorizer
}

func NewGroupsUsecase(groupsRepo *repository.GroupRepository, authz *Authorizer) *GroupsUsecase {
	return &GroupsUsecase{groupsRepo: groupsRepo, authz: authz}
}

func (uc *GroupsUsecase) Create(callerID uuid.UUID, req CreateGroupRequest) (*GroupDTO, error) {
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

	group := &domain.Group{
		Name:        name,
		Description: req.Description,
	}
	if err := uc.groupsRepo.Create(group); err != nil {
		return nil, err
	}
	return toGroupDTO(group), nil
}

func (uc *GroupsUsecase) GetByID(_ uuid.UUID, groupID uuid.UUID) (*GroupDTO, error) {
	group, err := uc.groupsRepo.GetByID(groupID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return toGroupDTO(group), nil
}

func (uc *GroupsUsecase) List(_ uuid.UUID) ([]*GroupDTO, error) {
	groups, err := uc.groupsRepo.List()
	if err != nil {
		return nil, err
	}
	out := make([]*GroupDTO, 0, len(groups))
	for _, g := range groups {
		out = append(out, toGroupDTO(g))
	}
	return out, nil
}

func (uc *GroupsUsecase) Update(callerID uuid.UUID, groupID uuid.UUID, req UpdateGroupRequest) (*GroupDTO, error) {
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

	if err := uc.groupsRepo.UpdateFields(groupID.String(), fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	group, err := uc.groupsRepo.GetByID(groupID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return toGroupDTO(group), nil
}

func (uc *GroupsUsecase) SoftDelete(callerID uuid.UUID, groupID uuid.UUID) error {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrForbidden
	}
	if err := uc.groupsRepo.SoftDelete(groupID.String()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func toGroupDTO(g *domain.Group) *GroupDTO {
	var ownerID *string
	if g.OwnerID != nil {
		s := g.OwnerID.String()
		ownerID = &s
	}
	return &GroupDTO{
		ID:          g.ID.String(),
		Name:        g.Name,
		Description: g.Description,
		OwnerID:     ownerID,
	}
}

