package product

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetProducts()(*[]Product, error) {
	rows, err := s.db.Query("SELECT id, name, description, price, stock_quanity, image_url, category_id, created_at, updated_at from products")
	if err != nil {
		return nil, err
	}
	products, err := s.ScanRowsIntoProduct(rows)
	if err != nil {
		return nil, err
	}
	return products, nil

}

func (s *Store) GetProductByID(id uuid.UUID)(*Product, error) {
	row := s.db.QueryRow("SELECT id, name, description, price, stock_quanity, image_url, category_id, created_at, updated_at FROM products WHERE id = $1", id)
	product, err := scanRowIntoProduct(row)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func scanRowIntoProduct(row *sql.Row) (*Product, error) {
	product := new(Product)
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.Image,
		&product.CategoryId,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return product, nil
}

func (s *Store) ScanRowsIntoProduct(rows *sql.Rows) (*[]Product, error) {
	
	products := make([]Product, 0)

	for rows.Next() {
		product := Product{}

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Image,
			&product.CategoryId,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return &products, nil
}