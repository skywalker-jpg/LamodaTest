package storage

import (
	"LamodaTest/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type WarehouseProductRepo struct {
	db *sql.DB
}

func NewWarehouseProductRepo(db *sql.DB) *WarehouseProductRepo {
	return &WarehouseProductRepo{
		db: db,
	}
}

func (r *WarehouseProductRepo) CreateWP(ctx context.Context, wp models.WarehouseProduct) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	insertQuery := squirrel.Insert("warehouse_product").
		Columns("warehouse_id", "product_id", "quantity", "reserved_quantity").
		Values(wp.WarehouseID, wp.ProductID, wp.Quantity, wp.ReservedQuantity).
		Suffix("RETURNING id").
		RunWith(tx).PlaceholderFormat(squirrel.Dollar)

	query, args, err := insertQuery.ToSql()
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert warehouse product: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return id, nil
}

func (r *WarehouseProductRepo) GetWP(ctx context.Context, filter models.GetWarehouseProductFilter) ([]*models.WarehouseProduct, error) {
	var warehouseProducts []*models.WarehouseProduct

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	queryBuilder := squirrel.Select("id", "warehouse_id", "product_id", "quantity", "reserved_quantity").From("warehouse_product").PlaceholderFormat(squirrel.Dollar)
	if len(filter.IDs) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": filter.IDs})
	}
	if filter.WarehouseID != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"warehouse_id": filter.WarehouseID})
	}
	if filter.ProductID != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"product_id": filter.ProductID})
	}

	queryBuilder = queryBuilder.RunWith(tx)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Println(query, args)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var wp models.WarehouseProduct
		if err := rows.Scan(&wp.ID, &wp.WarehouseID, &wp.ProductID, &wp.Quantity, &wp.ReservedQuantity); err != nil {
			return nil, fmt.Errorf("failed to scan warehouse products: %v", err)
		}
		warehouseProducts = append(warehouseProducts, &wp)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return warehouseProducts, nil
}

func (r *WarehouseProductRepo) GetWPByProductCode(ctx context.Context, filter models.GetWPByProductCodeFilter) (*models.WarehouseProduct, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	queryBuilder := squirrel.Select("wp.*").From("warehouse_product wp").
		Join("products p ON wp.product_id = p.id").
		Where(squirrel.Eq{"p.code": filter.ProductCode, "wp.warehouse_id": filter.WarehouseID}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)

	if filter.WarehouseID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"warehouse_id": *filter.WarehouseID})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var wp models.WarehouseProduct
	err = tx.QueryRowContext(ctx, query, args...).Scan(&wp.ID, &wp.WarehouseID, &wp.ProductID, &wp.Quantity, &wp.ReservedQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse product by product code: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &wp, nil
}

func (r *WarehouseProductRepo) UpdateWP(ctx context.Context, input *models.UpdateWarehouseProductInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	updateBuilder := squirrel.Update("warehouse_product").Where(squirrel.Eq{"id": input.ID}).RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	if input.WarehouseID != nil {
		updateBuilder = updateBuilder.Set("warehouse_id", *input.WarehouseID)
	}
	if input.ProductID != nil {
		updateBuilder = updateBuilder.Set("product_id", *input.ProductID)
	}
	if input.Quantity != nil {
		updateBuilder = updateBuilder.Set("quantity", *input.Quantity)
	}
	if input.ReservedQuantity != nil {
		updateBuilder = updateBuilder.Set("reserved_quantity", *input.ReservedQuantity)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update warehouse product: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *WarehouseProductRepo) UpdateWPBatch(ctx context.Context, inputs []models.UpdateWarehouseProductInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	for _, input := range inputs {
		updateBuilder := squirrel.Update("warehouse_product").Where(squirrel.Eq{"id": input.ID}).RunWith(tx).PlaceholderFormat(squirrel.Dollar)
		if input.WarehouseID != nil {
			updateBuilder = updateBuilder.Set("warehouse_id", *input.WarehouseID)
		}
		if input.ProductID != nil {
			updateBuilder = updateBuilder.Set("product_id", *input.ProductID)
		}
		if input.Quantity != nil {
			updateBuilder = updateBuilder.Set("quantity", *input.Quantity)
		}
		if input.ReservedQuantity != nil {
			updateBuilder = updateBuilder.Set("reserved_quantity", *input.ReservedQuantity)
		}

		query, args, err := updateBuilder.ToSql()
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update warehouse product: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *WarehouseProductRepo) DeleteWP(ctx context.Context, input models.DeleteWarehouseProductInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	deleteQuery := squirrel.Delete("warehouse_product").RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	if len(input.IDs) > 0 {
		deleteQuery = deleteQuery.Where(squirrel.Eq{"id": input.IDs})
	}
	if len(input.WarehouseID) > 0 {
		deleteQuery = deleteQuery.Where(squirrel.Eq{"warehouse_id": input.WarehouseID})
	}
	if len(input.ProductID) > 0 {
		deleteQuery = deleteQuery.Where(squirrel.Eq{"product_id": input.ProductID})
	}

	query, args, err := deleteQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse product: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
