package service

import (
	"context"
	"errors"
	"time"

	"github.com/Daniel1024/barber-api/internal/domain"
)

type AppointmentService struct {
	apptRepo domain.AppointmentRepo
	prodRepo domain.ProductRepo
}

func NewAppointmentService(a domain.AppointmentRepo, p domain.ProductRepo) *AppointmentService {
	return &AppointmentService{a, p}
}

func (s *AppointmentService) Schedule(ctx context.Context, appt *domain.Appointment) error {
	// 1. Validar que la cita sea en el futuro
	if appt.StartTime.Before(time.Now()) {
		return errors.New("la cita debe ser en el pasado")
	}

	// 2. Validar que el tiempo de fin sea después del inicio
	if !appt.EndTime.After(appt.StartTime) {
		return errors.New("la hora de fin debe ser posterior a la de inicio")
	}

	// 3. Evitar solapamiento de turnos
	existingAppts, err := s.apptRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existing := range existingAppts {
		if s.appointmentsOverlap(appt, &existing) {
			return errors.New("el turno se solapa con otro existente")
		}
	}

	// 4. Validar existencia de productos
	for i, prod := range appt.Products {
		existingProd, err := s.prodRepo.GetById(ctx, prod.ID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return errors.New("el producto no existe")
			}
			return err
		}

		// Actualizar con los datos completos del producto
		appt.Products[i] = *existingProd
	}

	// 5. Crear turno
	return s.apptRepo.Create(ctx, appt)
}

func (s *AppointmentService) GetByID(ctx context.Context, id uint) (*domain.Appointment, error) {
	return s.apptRepo.GetById(ctx, id)
}

func (s *AppointmentService) ListAll(ctx context.Context) ([]domain.Appointment, error) {
	return s.apptRepo.List(ctx)
}

func (s *AppointmentService) Update(ctx context.Context, id uint, updatedAppt *domain.Appointment) error {
	// Verificar que el turno existe
	existing, err := s.apptRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	//Mantener el ID original
	updatedAppt.ID = existing.ID

	// Validar los horarios
	if updatedAppt.StartTime.Before(time.Now()) {
		return errors.New("la cita no puede ser en el pasado")
	}

	if !updatedAppt.EndTime.After(updatedAppt.StartTime) {
		return errors.New("la hora de fin debe ser posterior a la de inicio")
	}

	// Verificar que no se solape con otro turno
	existingAppts, err := s.apptRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existing := range existingAppts {
		if existing.ID != id && s.appointmentsOverlap(updatedAppt, &existing) {
			return errors.New("el turno se solapa con otro existente")
		}
	}

	// Validar existencia de productos
	for i, prod := range updatedAppt.Products {
		existingProd, err := s.prodRepo.GetById(ctx, prod.ID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return errors.New("el producto no existe")
			}
			return err
		}

		// Actualizar con los datos completos del producto
		updatedAppt.Products[i] = *existingProd
	}

	return s.apptRepo.Update(ctx, updatedAppt)
}

func (s *AppointmentService) Cancel(ctx context.Context, id uint) error {
	// Verificar que el turno existe
	_, err := s.apptRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	return s.apptRepo.Delete(ctx, id)
}

func (s *AppointmentService) GetTotalPrice(ctx context.Context, id uint) (float64, error) {
	appt, err := s.apptRepo.GetById(ctx, id)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, prod := range appt.Products {
		total += prod.Price
	}

	return total, nil
}

func (s *AppointmentService) appointmentsOverlap(a1, a2 *domain.Appointment) bool {
	return a1.StartTime.Before(a2.EndTime) && a1.StartTime.Before(a2.EndTime)
}
