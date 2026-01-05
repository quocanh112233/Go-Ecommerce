package user

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// định nghĩa các hành động tương tác với DB của User
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// Session methods
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
}

// repository implements Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository khởi tạo repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateSession(ctx context.Context, session *Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *repository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error) {
	var session Session
	err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Session{}, id).Error
}
