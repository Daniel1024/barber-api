package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Daniel1024/barber-api/internal/domain"
	"github.com/Daniel1024/barber-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AppointmentHandle struct {
	svc *service.AppointmentService
}

func NewAppointmentHandle(svc *service.AppointmentService) *AppointmentHandle {
	return &AppointmentHandle{svc}
}

type CreateAppointmentRequest struct {
	ClientName string `json:"client_name" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
	Products   []uint `json:"products"`
}

type UpdateAppointmentRequest struct {
	ClientName string `json:"client_name" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
	Products   []uint `json:"products"`
}

func (h *AppointmentHandle) Create(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parseo de fechas
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha incorrecto para start_time, use RFC3339"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha incorrecto para end_time, use RFC3339"})
	}

	// Mapear IDs a domain.product
	products := make([]domain.Product, len(req.Products))
	for i, id := range req.Products {
		products[i] = domain.Product{ID: id}
	}

	appt := &domain.Appointment{
		ClientName: req.ClientName,
		StartTime:  startTime,
		EndTime:    endTime,
		Products:   products,
	}

	if err := h.svc.Schedule(c.Request.Context(), appt); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appt)
}

func (h *AppointmentHandle) List(c *gin.Context) {
	appoints, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments": appoints,
		"total":        len(appoints),
	})
}

func (h *AppointmentHandle) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	appt, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "turno no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandle) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parseo de fechas
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha incorrecto para start_time, use RFC3339"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha incorrecto para end_time, use RFC3339"})
	}

	// Mapear IDs a domain.product
	products := make([]domain.Product, len(req.Products))
	for i, id := range req.Products {
		products[i] = domain.Product{ID: id}
	}

	appt := &domain.Appointment{
		ID:         uint(id),
		ClientName: req.ClientName,
		StartTime:  startTime,
		EndTime:    endTime,
		Products:   products,
	}

	if err := h.svc.Update(c.Request.Context(), uint(id), appt); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandle) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.svc.Cancel(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "turno no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "turno cancelado"})
}

func (h *AppointmentHandle) GetTotal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	total, err := h.svc.GetTotalPrice(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "turno no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
