package repositories

import (
	"context"
	"tikcitz-app/internals/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersRepository struct {
	db *pgxpool.Pool
}

func NewOrdersRepository(db *pgxpool.Pool) *OrdersRepository {
	return &OrdersRepository{db: db}
}

// repository create order
func (o *OrdersRepository) CreateOrder(ctx context.Context, order models.OrdersStruct, IdInt int) (pgconn.CommandTag, error) {
	query := "INSERT INTO orders(user_id, movie_id, total_price, full_name, email, phone_number, payment_methode_id, paid, date, time, cinema_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	values := []any{IdInt, order.Movie_id, order.Total_price, order.Full_name, order.Email, order.Phone_number, order.Payment_methode_id, order.Paid, order.Date, order.Time, order.Cinema_id}

	cmd, err := o.db.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}

