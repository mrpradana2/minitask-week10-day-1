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

// repository get order history user
func (o *OrdersRepository) GetOrderHistory(ctx context.Context, IdInt int) ([]models.OrdersStruct, error) {
	query := "SELECT o.user_id, o.total_price, o.phone_number, o.paid, pm.name, m.title, c.image_path FROM orders o JOIN payment_methode pm ON o.payment_methode_id = pm.id JOIN movies m ON m.id = o.movie_id JOIN cinemas c ON o.cinema_id = c.id WHERE o.user_id = $1"
	values := []any{IdInt}
	rows, err := o.db.Query(ctx, query, values...)
	if err != nil {
		return nil, err 	
	} 
	
	defer rows.Close()

	var result []models.OrdersStruct
	for rows.Next() {
		var order models.OrdersStruct
		err := rows.Scan(&order.User_id, &order.Total_price, &order.Phone_number, &order.Paid, &order.Payment_methode, &order.Title, &order.Cinema_path)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}

	return result, nil
}
