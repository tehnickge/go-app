package usecase

import (
	"errors"
	"strings"
	"time"

	"backend-app/internal/domain"
	"backend-app/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	usersRepo *repository.UserRepository
	rolesRepo *repository.RoleRepository
	urRepo    *repository.UsersRolesRepository

	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthUsecase(
	usersRepo *repository.UserRepository,
	rolesRepo *repository.RoleRepository,
	urRepo *repository.UsersRolesRepository,
	jwtSecret string,
	jwtTTL time.Duration,
) *AuthUsecase {
	return &AuthUsecase{
		usersRepo: usersRepo,
		rolesRepo: rolesRepo,
		urRepo:    urRepo,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

type RegisterRequest struct {
	Login    string
	Email    string
	Password string
}

type RegisterResponse struct {
	UserID string
	Token  string
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
}

func (uc *AuthUsecase) Register(req RegisterRequest) (*RegisterResponse, error) {
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
	}

	if err := uc.usersRepo.Create(user); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}

	// Bootstrap admin role:
	// - ensure "admin" role exists
	// - assign it to the first registered user if no active admin assignments exist
	adminRole, err := uc.rolesRepo.GetByName("admin")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = uc.rolesRepo.Create(&domain.Role{Name: "admin", Description: "global admin"})
			adminRole, err = uc.rolesRepo.GetByName("admin")
		}
	}
	if err == nil && adminRole != nil {
		hasAdmin, err := uc.urRepo.HasAnyActiveAssignmentByRoleID(adminRole.ID)
		if err != nil {
			return nil, err
		}
		if !hasAdmin {
			_ = uc.urRepo.AssignRoleToUser(user.ID, adminRole.ID)
		}
	}

	token, err := uc.issueToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		UserID: user.ID.String(),
		Token:  token,
	}, nil
}

func (uc *AuthUsecase) Login(req LoginRequest) (*LoginResponse, error) {
	email := strings.TrimSpace(req.Email)
	if email == "" || req.Password == "" {
		return nil, errors.New("invalid input")
	}

	user, err := uc.usersRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrUnauthorized
	}

	token, err := uc.issueToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token}, nil
}

func (uc *AuthUsecase) issueToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(uc.jwtTTL).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(uc.jwtSecret))
}

func isUniqueViolation(err error) bool {
	// Best-effort mapping. Postgres error text typically contains "duplicate key".
	return strings.Contains(err.Error(), "duplicate key")
}

