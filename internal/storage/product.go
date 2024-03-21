package storage

import (
	"LamodaTest/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) CreateProduct(ctx context.Context, p models.Product) (int, error) {
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

	insertQuery := squirrel.Insert("products").
		Columns("name", "size", "code").
		Values(p.Name, p.Size, p.Code).
		Suffix("RETURNING id").
		RunWith(tx).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := insertQuery.ToSql()
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRowContext(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert product: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return id, nil
}

func (r *ProductRepo) GetProducts(ctx context.Context, filter models.GetProductsFilter) ([]*models.Product, error) {
	var products []*models.Product

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

	queryBuilder := squirrel.Select("id", "name", "size", "code").From("products").RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	if len(filter.IDs) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": filter.IDs})
	}
	if len(filter.Codes) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"code": filter.Codes})
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
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Size, &product.Code); err != nil {
			return nil, fmt.Errorf("failed to scan products: %v", err)
		}
		products = append(products, &product)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return products, nil
}

func (r *ProductRepo) UpdateProduct(ctx context.Context, input *models.UpdateProductInput) error {
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

	updateBuilder := squirrel.Update("products").Where(squirrel.Eq{"id": input.ID}).RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", *input.Name)
	}
	if input.Size != nil {
		updateBuilder = updateBuilder.Set("size", *input.Size)
	}
	if input.Code != nil {
		updateBuilder = updateBuilder.Set("code", *input.Code)
	}
	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot update product: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, input models.DeleteProductInput) error {
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

	deleteQuery := squirrel.Delete("products").Where(squirrel.Eq{"id": input.ID}).RunWith(tx).PlaceholderFormat(squirrel.Dollar)
	query, args, err := deleteQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
