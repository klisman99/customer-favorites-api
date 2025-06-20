package db

import (
	"app/internal/domain/model"
	"context"
	"database/sql"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

func (r *CustomerRepository) Create(c context.Context, customer model.Customer) (*model.Customer, error) {
	query := `
		INSERT INTO customers (id, name, email)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at
	`

	row := r.DB.QueryRowContext(c, query, customer.ID, customer.Name, customer.Email)

	var createdCustomer model.Customer
	if err := row.Scan(
		&createdCustomer.ID,
		&createdCustomer.Name,
		&createdCustomer.Email,
		&createdCustomer.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &createdCustomer, nil
}

func (r *CustomerRepository) FindByID(c context.Context, id string) (*model.Customer, error) {
	query := `
		SELECT id, name, email, created_at
		FROM customers
		WHERE id = $1
	`

	row := r.DB.QueryRowContext(c, query, id)

	var customer model.Customer
	if err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) FindAll(c context.Context) ([]model.Customer, error) {
	query := `
		SELECT id, name, email, created_at
		FROM customers
	`

	rows, err := r.DB.QueryContext(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		if err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.CreatedAt,
		); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func (r *CustomerRepository) Update(c context.Context, customer model.Customer) error {
	query := `
		UPDATE customers
		SET name = $1, email = $2, created_at = $3
		WHERE id = $4
	`

	_, err := r.DB.ExecContext(c, query,
		customer.Name,
		customer.Email,
		customer.CreatedAt,
		customer.ID,
	)
	return err
}

func (r *CustomerRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM customers WHERE id = $1"

	_, err := r.DB.ExecContext(c, query, id)
	return err
}

func (r *CustomerRepository) FindByEmail(c context.Context, email string) (*model.Customer, error) {
	query := `
		SELECT id, name, email, created_at
		FROM customers
		WHERE email = $1
	`

	row := r.DB.QueryRowContext(c, query, email)

	var customer model.Customer
	if err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}
