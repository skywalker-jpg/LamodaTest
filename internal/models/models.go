package models

// Warehouse represents model for warehouses table
type Warehouse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Availability bool   `json:"availability"`
}

type GetWarehousesFilter struct {
	IDs []int `json:"ID,omitempty"`
}

type UpdateWarehouseInput struct {
	ID           int     `json:"id"`
	Name         *string `json:"name"`
	Availability *bool   `json:"availability"`
}

type DeleteWarehouseInput struct {
	ID int
}

// Product represents model for products table
type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Size string `json:"size"`
	Code string `json:"code"`
}

type GetProductsFilter struct {
	IDs   []int    `json:"IDs,omitempty"`
	Codes []string `json:"code,omitempty"`
}

type UpdateProductInput struct {
	ID int `json:"id"`

	Name *string `json:"name"`
	Size *string `json:"size"`
	Code *string `json:"code"`
}

type DeleteProductInput struct {
	ID int
}

// WarehouseProduct represents model for warehouse_product table
type WarehouseProduct struct {
	ID               int `json:"id"`
	WarehouseID      int `json:"warehouse_id"`
	ProductID        int `json:"product_id"`
	Quantity         int `json:"quantity"`
	ReservedQuantity int `json:"reserved_quantity"`
}

type GetWarehouseProductFilter struct {
	IDs         []int `json:"IDs,omitempty"`
	WarehouseID int   `json:"WarehouseID,omitempty"`
	ProductID   int   `json:"ProductID,omitempty"`
}

type GetWPByProductCodeFilter struct {
	WarehouseID *int    `json:"WarehouseID,omitempty"`
	ProductCode *string `json:"ProductID,omitempty"`
}

type UpdateWarehouseProductInput struct {
	ID int `json:"id"`

	WarehouseID      *int `json:"warehouse_id"`
	ProductID        *int `json:"product_id"`
	Quantity         *int `json:"quantity"`
	ReservedQuantity *int `json:"reserved_quantity"`
}

type DeleteWarehouseProductInput struct {
	IDs         []int `json:"IDs,omitempty"`
	WarehouseID []int `json:"WarehouseIDs,omitempty"`
	ProductID   []int `json:"ProductIDs,omitempty"`
}

type ReserveDTO struct {
	Reservations []Reserve `json:"reservations"`
	WarehouseID  int       `json:"warehouse_id"`
}

type Reserve struct {
	ProductID       int `json:"product_id"`
	ReserveQuantity int `json:"reserve_quantity"`
}
