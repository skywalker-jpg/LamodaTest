package main

import (
	"LamodaTest/internal/config"
	"LamodaTest/internal/database"
	"LamodaTest/internal/logger"
	"LamodaTest/internal/models"
	"LamodaTest/internal/storage"
	"LamodaTest/internal/web"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	config, err := config.NewConfig("config.yaml")
	if err != nil {
		return err
	}

	logger, err := logger.New(config.Logger)
	if err != nil {
		return err
	}

	db, err := database.Connection(config.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	err = database.MigrateUp(config.DB)
	if err != nil {
		logger.Error("Error creating migrations", err.Error())
		return err
	}

	ctx := context.TODO()

	st := storage.NewStorage(db)

	err = outputTestData(ctx, st)
	if err != nil {
		return err
	}

	server, err := web.New(config.Server, logger, st)
	if err != nil {
		return err
	}

	return server.Serve()
}

func outputTestData(ctx context.Context, st *storage.Storage) error {
	warehouses, err := st.GetWarehouses(ctx, models.GetWarehousesFilter{})
	if err != nil {
		return fmt.Errorf("error getting TestData %v", err)
	}
	products, err := st.GetProducts(ctx, models.GetProductsFilter{})
	if err != nil {
		return fmt.Errorf("error getting TestData %v", err)
	}
	warehouseProducts, err := st.GetWP(ctx, models.GetWarehouseProductFilter{})

	for _, warehouse := range warehouses {
		jsonData, err := json.MarshalIndent(warehouse, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d, Name: %s, Availability: %t\n%s\n", warehouse.ID, warehouse.Name, warehouse.Availability, string(jsonData))
	}
	for _, product := range products {
		jsonData, err := json.MarshalIndent(product, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d, Name: %s, Size: %s, Code: %s\n%s\n", product.ID, product.Name, product.Size, product.Code, string(jsonData))
	}
	for _, warehouseProduct := range warehouseProducts {
		jsonData, err := json.MarshalIndent(warehouseProduct, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d, WarehouseID %d, ProductId: %d, Quantity: %d, ReservedQuantity: %d\n%s\n", warehouseProduct.ID, warehouseProduct.WarehouseID, warehouseProduct.ProductID, warehouseProduct.Quantity, warehouseProduct.ReservedQuantity, string(jsonData))
	}

	return nil
}
