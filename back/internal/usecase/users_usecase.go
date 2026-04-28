package usecase

import (
	"errors"
	"strings"

	"backend-app/internal/domain"
	"backend-app/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID     string `json:"id"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type CreateUserRequest struct {
	Login    string
	Email    string
	Password string
}

type UpdateUserRequest struct {
	Login    *string
	Email    *string
	Password *string
	Status   *string
}

type UsersUsecase struct {
	usersRepo *repository.UserRepository
	authz     *Authorizer
}

func NewUsersUsecase(usersRepo *repository.UserRepository, authz *Authorizer) *UsersUsecase {
	return &UsersUsecase{usersRepo: usersRepo, authz: authz}
}

func (uc *UsersUsecase) Create(callerID uuid.UUID, req CreateUserRequest) (*UserDTO, error) {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	login := strings.TrimSpace(req.Login)
	email := strings.TrimSpace(req.Email)
	pass := req.Password
	if login == "" || email == "" || pass == "" {
		return nil, errors.New("invalid input")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Login:    login,
		Email:    email,
		Password: string(hashed),
		Status:   "active",
		// timestamps set by gorm
	}

	if err := uc.usersRepo.Create(user); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}

	return toUserDTO(user), nil
}

func (uc *UsersUsecase) GetByID(callerID uuid.UUID, userID uuid.UUID) (*UserDTO, error) {
	if callerID != userID {
		isAdmin, err := uc.authz.IsAdmin(callerID)
		if err != nil {
			return nil, err
		}
		if !isAdmin {
			return nil, ErrForbidden
		}
	}

	user, err := uc.usersRepo.GetByID(userID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return toUserDTO(user), nil
}

func (uc *UsersUsecase) List(callerID uuid.UUID) ([]*UserDTO, error) {
	isAdmin, err := uc.authz.IsAdmin(callerID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	users, err := uc.usersRepo.List()
	if err != nil {
		return nil, err
	}

	out := make([]*UserDTO, 0, len(users))
	for _, u := range users {
		out = append(out, toUserDTO(u))
	}
	return out, nil
}

func (uc *UsersUsecase) Update(callerID uuid.UUID, userID uuid.UUID, req UpdateUserRequest) (*UserDTO, error) {
	if callerID != userID {
		isAdmin, err := uc.authz.IsAdmin(callerID)
		if err != nil {
			return nil, err
		}
		if !isAdmin {
			return nil, ErrForbidden
		}
	}

	fields := map[string]interface{}{}
	if req.Login != nil {
		login := strings.TrimSpace(*req.Login)
		if login == "" {
			return nil, errors.New("invalid input")
		}
		fields["login"] = login
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" {
			return nil, errors.New("invalid input")
		}
		fields["email"] = email
	}
	if req.Password != nil {
		pass := *req.Password
		if pass == "" {
			return nil, errors.New("invalid input")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		fields["password"] = string(hashed)
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}

	if len(fields) == 0 {
		return nil, errors.New("nothing to update")
	}

	if err := uc.usersRepo.UpdateFields(userID.String(), fields); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user, err := uc.usersRepo.GetByID(userID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return toUserDTO(user), nil
}

func (uc *UsersUsecase) SoftDelete(callerID uuid.UUID, userID uuid.UUID) error {
	if callerID != userID {
		isAdmin, err := uc.authz.IsAdmin(callerID)
		if err != nil {
			return err
		}
		if !isAdmin {
			return ErrForbidden
		}
	}

	if err := uc.usersRepo.SoftDelete(userID.String()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func toUserDTO(u *domain.User) *UserDTO {
	return &UserDTO{
		ID:     u.ID.String(),
		Login:  u.Login,
		Email:  u.Email,
		Status: u.Status,
	}
}

