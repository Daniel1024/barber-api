package repository

import (
	"context"
	"errors"

	"github.com/Daniel1024/barber-api/internal/domain"
	"gorm.io/gorm"
)

type GormProductRepo struct {
	db *gorm.DB
}

func NewGormProduct(db *gorm.DB) domain.ProductRepo {
	return &GormProductRepo{db}
}

func (r *GormProductRepo) Create(ctx context.Context, appt *domain.Product) error {
	return r.db.WithContext(ctx).Create(appt).Error
}

func (r *GormProductRepo) GetById(ctx context.Context, id uint) (*domain.Product, error) {
	var appt domain.Product

	if err := r.db.WithContext(ctx).Preload("Products").First(&appt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &appt, nil
}

func (r *GormProductRepo) List(ctx context.Context) ([]domain.Product, error) {
	var appts []domain.Product

	if err := r.db.WithContext(ctx).Preload("Products").Find(&appts).Error; err != nil {
		return nil, err
	}

	return appts, nil
}

func (r *GormProductRepo) Update(ctx context.Context, appt *domain.Product) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(appt).Error
}

func (r *GormProductRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}
