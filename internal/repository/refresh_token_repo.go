package repository

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) interfaces.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.Where("token = ? AND revoked = false", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepository) Revoke(token string) error {
	return r.db.Model(&domain.RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}

func (r *refreshTokenRepository) RevokeAllForUser(userID string) error {
	return r.db.Model(&domain.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
