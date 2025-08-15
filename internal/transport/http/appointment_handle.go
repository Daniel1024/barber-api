package http

import "github.com/Daniel1024/barber-api/internal/service"

type AppointmentHandle struct {
	svc *service.AppointmentService
}

func NewAppointmentHandle(svc *service.AppointmentService) *AppointmentHandle {
	return &AppointmentHandle{svc}
}

type CreateAppointmentRequest struct {
	ClientName string `json:"client_name" binding:"required"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Products   []uint `json:"products"`
}
