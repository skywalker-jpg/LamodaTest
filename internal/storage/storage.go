package storage

import (
	"LamodaTest/internal/models"
	"context"
	"database/sql"
)

type WarehouseStorage interface {
	CreateWarehouse(ctx context.Context, wh models.Warehouse) (int, error)
	GetWarehouses(ctx context.Context, filter models.GetWarehousesFilter) ([]*models.Warehouse, error)
	UpdateWarehouse(ctx context.Context, input *models.UpdateWarehouseInput) error
	DeleteWarehouse(ctx context.Context, wh models.DeleteWarehouseInput) error
}

type ProductStorage interface {
	CreateProduct(ctx context.Context, p models.Product) (int, error)
	GetProducts(ctx context.Context, filter models.GetProductsFilter) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, input *models.UpdateProductInput) error
	DeleteProduct(ctx context.Context, input models.DeleteProductInput) error
}

type WarehouseProductStorage interface {
	CreateWP(ctx context.Context, wp models.WarehouseProduct) (int, error)
	GetWP(ctx context.Context, filter models.GetWarehouseProductFilter) ([]*models.WarehouseProduct, error)
	GetWPByProductCode(ctx context.Context, filter models.GetWPByProductCodeFilter) (*models.WarehouseProduct, error)
	UpdateWP(ctx context.Context, input *models.UpdateWarehouseProductInput) error
	UpdateWPBatch(ctx context.Context, inputs []models.UpdateWarehouseProductInput) error
	DeleteWP(ctx context.Context, input models.DeleteWarehouseProductInput) error
}

type Storage struct {
	WarehouseStorage
	ProductStorage
	WarehouseProductStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		WarehouseStorage:        NewWarehouseRepo(db),
		ProductStorage:          NewProductRepo(db),
		WarehouseProductStorage: NewWarehouseProductRepo(db),
	}
}
