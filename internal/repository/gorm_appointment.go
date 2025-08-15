package repository

import (
	"context"
	"errors"

	"github.com/Daniel1024/barber-api/internal/domain"
	"gorm.io/gorm"
)

type GormAppointmentRepo struct {
	db *gorm.DB
}

func NewGormAppointment(db *gorm.DB) domain.AppointmentRepo {
	return &GormAppointmentRepo{db}
}

func (r *GormAppointmentRepo) Create(ctx context.Context, appt *domain.Appointment) error {
	return r.db.WithContext(ctx).Create(appt).Error
}

func (r *GormAppointmentRepo) GetById(ctx context.Context, id uint) (*domain.Appointment, error) {
	var appt domain.Appointment

	if err := r.db.WithContext(ctx).Preload("Products").First(&appt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &appt, nil
}

func (r *GormAppointmentRepo) List(ctx context.Context) ([]domain.Appointment, error) {
	var appts []domain.Appointment

	if err := r.db.WithContext(ctx).Preload("Products").Find(&appts).Error; err != nil {
		return nil, err
	}

	return appts, nil
}

func (r *GormAppointmentRepo) Update(ctx context.Context, appt *domain.Appointment) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(appt).Error
}

func (r *GormAppointmentRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Appointment{}, id).Error
}
