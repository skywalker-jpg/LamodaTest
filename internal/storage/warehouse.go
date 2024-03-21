package storage

import (
	"LamodaTest/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type WarehouseRepo struct {
	db *sql.DB
}

func NewWarehouseRepo(db *sql.DB) *WarehouseRepo {
	return &WarehouseRepo{
		db: db,
	}
}

func (r *WarehouseRepo) CreateWarehouse(ctx context.Context, warehouse models.Warehouse) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	insertQuery := squirrel.Insert("warehouses").
		Columns("name", "availability").
		Values(warehouse.Name, warehouse.Availability).
		Suffix("RETURNING id").
		RunWith(tx).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := insertQuery.ToSql()
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRowContext(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *WarehouseRepo) GetWarehouses(ctx context.Context, filter models.GetWarehousesFilter) ([]*models.Warehouse, error) {
	var warehouses []*models.Warehouse

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	queryBuilder := squirrel.Select("id", "name", "availability").From("warehouses").RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	if len(filter.IDs) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": filter.IDs})
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var warehouse models.Warehouse
		if err := rows.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Availability); err != nil {
			return nil, fmt.Errorf("failed to scan warehouses: %v", err)
		}
		warehouses = append(warehouses, &warehouse)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return warehouses, nil
}

func (r *WarehouseRepo) UpdateWarehouse(ctx context.Context, input *models.UpdateWarehouseInput) error {
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

	updateBuilder := squirrel.Update("warehouses").Where(squirrel.Eq{"id": input.ID}).RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", *input.Name)
	}
	if input.Availability != nil {
		updateBuilder = updateBuilder.Set("availability", *input.Availability)
	}
	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot update warehouse: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *WarehouseRepo) DeleteWarehouse(ctx context.Context, input models.DeleteWarehouseInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				fmt.Printf("Rollback error: %v\n", rollbackErr)
			}
		}
	}()

	deleteQuery := squirrel.Delete("warehouses").Where(squirrel.Eq{"id": input.ID}).RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	query, args, err := deleteQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
