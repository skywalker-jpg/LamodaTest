package web

import (
	"LamodaTest/internal/models"
	"context"
	"fmt"
	"github.com/labstack/echo"
	"log/slog"
	"net/http"
)

type ReserveDTO struct {
	Reservations []Reserve `json:"reservations"`
}

type Reserve struct {
	Code        string `json:"code"`
	Quantity    int    `json:"quantity"`
	WarehouseID int    `json:"warehouse_id"`
}

type ReleaseDTO struct {
	Releases []Release `json:"releases"`
}

type Release struct {
	Code        string `json:"code"`
	Quantity    int    `json:"quantity"`
	WarehouseID int    `json:"warehouse_id"`
}

type WarehouseProductsDTO struct {
	WarehouseID int `json:"warehouse_id"`
}

type BlockWarehouseDTO struct {
	WarehouseID int `json:"warehouse_id"`
}

func (s *Server) RegisterHandlers() {
	app := s.app

	apiGroup := app.Group("/api/v1")
	apiGroup.POST("/reserve", s.ReserveProductHandler)
	apiGroup.POST("/release", s.ReleaseProductHandler)

	apiGroup.POST("/products", s.GetWarehouseHandler)
	apiGroup.POST("/block", s.BlockWarehouseHandler)
	apiGroup.POST("/unblock", s.UnblockWarehouseHandler)

	app.GET("/*", s.NotFound)
}

func (s *Server) ReserveProductHandler(c echo.Context) error {
	requestID := c.Get("requestID").(string)
	var reserveData ReserveDTO
	if err := c.Bind(&reserveData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if len(reserveData.Reservations) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Empty request"})
	}

	updates := make([]models.UpdateWarehouseProductInput, len(reserveData.Reservations))

	for i, reservation := range reserveData.Reservations {
		wp, err := s.Storage.GetWPByProductCode(context.TODO(), models.GetWPByProductCodeFilter{WarehouseID: &reservation.WarehouseID, ProductCode: &reservation.Code})
		if err != nil {
			s.logger.Error("Server", slog.String("requestID", requestID),
				slog.String("error", fmt.Sprintf("Unable to get stored products in warehouse: %v", err.Error())))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to get stored products in warehouse: %v", err.Error())})
		}
		if wp == nil {
			s.logger.Info("Server", slog.String("requestID", requestID),
				slog.String("error", "No products in warehouse"))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No products in warehouse"})
		}

		reserved := wp.ReservedQuantity + reservation.Quantity
		if reserved > wp.Quantity {
			s.logger.Info("Server", slog.String("requestID", requestID),
				slog.String("error", "Can't reserve more than have"))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Can't reserve more than have"})
		}

		updates[i] = models.UpdateWarehouseProductInput{WarehouseID: &wp.WarehouseID, ProductID: &wp.ProductID, Quantity: &wp.Quantity, ReservedQuantity: &reserved}
	}
	err := s.Storage.UpdateWPBatch(context.TODO(), updates)
	if err != nil {
		s.logger.Error("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("Unable to update warehouse_product records: %v", err.Error())))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to update warehouse_product records: %v", err.Error())})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Reserved": "OK"})

}

func (s *Server) ReleaseProductHandler(c echo.Context) error {
	requestID := c.Get("requestID").(string)
	var releaseData ReleaseDTO
	if err := c.Bind(&releaseData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if len(releaseData.Releases) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Empty request"})
	}

	updates := make([]models.UpdateWarehouseProductInput, len(releaseData.Releases))

	for i, reservation := range releaseData.Releases {
		wp, err := s.Storage.GetWPByProductCode(context.TODO(), models.GetWPByProductCodeFilter{WarehouseID: &reservation.WarehouseID, ProductCode: &reservation.Code})
		if err != nil {
			s.logger.Error("Server", slog.String("requestID", requestID),
				slog.String("error", fmt.Sprintf("Unable to get stored products in warehouse: %v", err.Error())))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to get stored products in warehouse: %v", err.Error())})
		}
		if wp == nil {
			s.logger.Info("Server", slog.String("requestID", requestID),
				slog.String("error", "No products in warehouse"))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No products in warehouse"})
		}

		released := wp.ReservedQuantity - reservation.Quantity
		if released < 0 {
			s.logger.Info("Server", slog.String("requestID", requestID),
				slog.String("error", "Can't release more than have"))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Can't release more than have"})
		}

		updates[i] = models.UpdateWarehouseProductInput{WarehouseID: &wp.WarehouseID, ProductID: &wp.ProductID, Quantity: &wp.Quantity, ReservedQuantity: &released}
	}
	err := s.Storage.UpdateWPBatch(context.TODO(), updates)
	if err != nil {
		s.logger.Error("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("Unable to update warehouse_product records: %v", err.Error())))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to update warehouse_product records: %v", err.Error())})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Released": "OK"})

}

func (s *Server) GetWarehouseHandler(c echo.Context) error {
	requestID := c.Get("requestID").(string)
	var warehouseProducts WarehouseProductsDTO
	if err := c.Bind(&warehouseProducts); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	wp, err := s.Storage.GetWP(context.TODO(), models.GetWarehouseProductFilter{WarehouseID: warehouseProducts.WarehouseID})
	if err != nil {
		fmt.Println(&warehouseProducts.WarehouseID)
		s.logger.Error("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("Unable to get warehouse products: %v", err.Error())))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to get warehouse products: %v", err.Error())})
	}

	if wp == nil {
		s.logger.Info("Server", slog.String("requestID", requestID),
			slog.String("error", "No products in warehouse"))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No products in warehouse"})
	}

	return c.JSON(http.StatusOK, wp)

}

func (s *Server) BlockWarehouseHandler(c echo.Context) error {
	requestID := c.Get("requestID").(string)
	False := false
	var warehouse BlockWarehouseDTO
	if err := c.Bind(&warehouse); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	wps, err := s.Storage.GetWarehouses(context.TODO(), models.GetWarehousesFilter{IDs: []int{warehouse.WarehouseID}})
	if err != nil {
		s.logger.Error("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("Unable to get warehouse: %v", err.Error())))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to get warehouse: %v", err.Error())})
	}
	if wps == nil {
		s.logger.Info("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("No warehouse with ID: %d", warehouse.WarehouseID)))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("No warehouse with ID: %d", warehouse.WarehouseID)})
	}
	err = s.Storage.UpdateWarehouse(context.TODO(), &models.UpdateWarehouseInput{ID: wps[0].ID, Name: &wps[0].Name, Availability: &False})
	if err != nil {
		s.logger.Error("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("Unable to update warehouse: %v", err.Error())))
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": fmt.Sprintf("Unable to update warehouse: %v", err.Error())})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Blocked": "OK",
	})
}

func (s *Server) UnblockWarehouseHandler(c echo.Context) error {
	requestID := c.Get("requestID").(string)
	True := true
	var warehouse BlockWarehouseDTO
	if err := c.Bind(&warehouse); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	wps, err := s.Storage.GetWarehouses(context.TODO(), models.GetWarehousesFilter{IDs: []int{warehouse.WarehouseID}})
	if err != nil {
		s.logger.Error("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("Unable to get warehouse: %v", err.Error())))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Unable to get warehouse: %v", err.Error())})
	}
	if wps == nil {
		s.logger.Info("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("No warehouse with ID: %d", warehouse.WarehouseID)))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("No warehouse with ID: %d", warehouse.WarehouseID)})
	}
	err = s.Storage.UpdateWarehouse(context.TODO(), &models.UpdateWarehouseInput{ID: wps[0].ID, Name: &wps[0].Name, Availability: &True})
	if err != nil {
		s.logger.Error("Server", slog.String("requestID", requestID),
			slog.String("error", fmt.Sprintf("Unable to update warehouse: %v", err.Error())))
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": fmt.Sprintf("Unable to update warehouse: %v", err.Error())})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Unblocked": "OK",
	})
}

func (s *Server) NotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, "Page not found")
}
