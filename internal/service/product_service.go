package service

import (
	"context"
	"errors"

	"github.com/Daniel1024/barber-api/internal/domain"
)

type ProductService struct {
	productRepo domain.ProductRepo
}

func NewProductService(p domain.ProductRepo) *ProductService {
	return &ProductService{p}
}

func (s *ProductService) Create(ctx context.Context, product *domain.Product) error {
	// Validaciones básicas
	if product.Name == "" {
		return errors.New("el nombre del producto no puede estar vació")
	}

	if product.Price <= 0 {
		return errors.New("el precio del producto debe ser mayor a 0")
	}

	// Verificar que no existe un producto con el mismo nombre
	existing, err := s.productRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existingProd := range existing {
		if existingProd.Name == product.Name {
			return errors.New("ya existe un producto con ese nombre")
		}
	}

	return s.productRepo.Create(ctx, product)
}

func (s *ProductService) ListAll(ctx context.Context) ([]domain.Product, error) {
	return s.productRepo.List(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	return s.productRepo.GetById(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, id uint, updatedProduct *domain.Product) error {
	// Validaciones básicas
	if updatedProduct.Name == "" {
		return errors.New("el nombre del producto no puede estar vació")
	}

	if updatedProduct.Price <= 0 {
		return errors.New("el precio del producto debe ser mayor a 0")
	}

	// Verificar que el turno existe
	existing, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	// Verificar que no existe un producto con el mismo nombre
	existingProducts, err := s.productRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, existingProd := range existingProducts {
		if existingProd.Name == updatedProduct.Name && existingProd.ID != id {
			return errors.New("ya existe un producto con ese nombre")
		}
	}

	//Mantener el ID original
	updatedProduct.ID = existing.ID

	return s.productRepo.Update(ctx, updatedProduct)
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	// Verificar que el turno existe
	_, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(ctx, id)
}
