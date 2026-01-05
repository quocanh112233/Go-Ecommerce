package user

import (
	"context"
	"time"

	"go-ecommerce/internal/config"
	"go-ecommerce/internal/shared/errors"
	"go-ecommerce/pkg/crypto"
	"go-ecommerce/pkg/token"

	"github.com/google/uuid"
)

// Service interface định nghĩa các method mà tầng Handler sẽ gọi
type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
}

// service struct implement interface trên
type service struct {
	repo Repository
	cfg  *config.Config
}

// NewService khởi tạo service
func NewService(repo Repository, cfg *config.Config) Service {
	return &service{repo: repo, cfg: cfg}
}

// Register thực hiện logic đăng ký
func (s *service) Register(ctx context.Context, req RegisterRequest) (*UserResponse, error) {
	// 1. Kiểm tra Email đã tồn tại chưa
	exists, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && exists != nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	// 2. Hash mật khẩu
	hashedPassword := crypto.HashPassword(req.Password)

	// 3. Tạo Entity User từ Request
	newUser := &User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Username:     req.Username,
		Phone:        req.Phone,
		Role:         RoleCustomer, // Mặc định là Customer
		IsActive:     true,
	}

	// 4. Lưu vào Database
	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// 5. Map từ Entity sang Response DTO để trả về
	return &UserResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Username:  newUser.Username,
		Phone:     newUser.Phone,
		Role:      string(newUser.Role),
		CreatedAt: newUser.CreatedAt,
	}, nil
}

// Login xử lý đăng nhập
func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// 1. Tìm user theo Email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 2. Kiểm tra password
	if !crypto.ComparePassword(user.PasswordHash, req.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	// 3. Generate Token
	accessToken, err := token.GenerateAccessToken(user.ID, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.AccessExpiration)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := token.GenerateRefreshToken()

	// 4. Lưu Session
	session := &Session{
		UserID:       user.ID,
		RefreshToken: refreshTokenStr,
		UserAgent:    "",
		ClientIP:     "",
		ExpiresAt:    time.Now().Add(s.cfg.JWT.RefreshExpiration),
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// 5. Trả về
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(s.cfg.JWT.AccessExpiration.Seconds()),
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Phone:     user.Phone,
			Role:      string(user.Role),
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// RefreshToken tạo access token mới từ refresh token
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// 1. Tìm Session
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 2. Kiểm tra hết hạn hoặc bị block
	if session.IsBlocked {
		return nil, errors.ErrInvalidCredentials
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.ErrInvalidCredentials
	}

	// 3. Lấy User info để tạo token mới
	user, err := s.repo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 4. Generate Access Token mới
	accessToken, err := token.GenerateAccessToken(user.ID, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.AccessExpiration)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Giữ nguyên refresh token cũ
		ExpiresIn:    int64(s.cfg.JWT.AccessExpiration.Seconds()),
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Phone:     user.Phone,
			Role:      string(user.Role),
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// Logout xóa session khi đăng xuất
func (s *service) Logout(ctx context.Context, refreshToken string) error {
	// Tìm Session
	session, err := s.repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.ErrInvalidCredentials
	}

	// Xóa Session
	return s.repo.DeleteSession(ctx, session.ID)
}

// GetProfile lấy thông tin user theo ID
func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Phone:     user.Phone,
		Role:      string(user.Role),
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}, nil
}
