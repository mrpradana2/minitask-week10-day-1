package repositories

import (
	"context"
	"tikcitz-app/internals/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersRepository struct {
	db *pgxpool.Pool
}

func NewOrdersRepository(db *pgxpool.Pool) *OrdersRepository {
	return &OrdersRepository{db: db}
}

// repository create order
func (o *OrdersRepository) CreateOrder(ctx context.Context, order models.OrdersStruct, IdInt int) error {

	// insert to table orders and returning id and cinema_id
	query := "INSERT INTO orders(user_id, movie_id, total_price, full_name, email, phone_number, payment_methode_id, paid, date, time, cinema_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, cinema_id" //
	values := []any{IdInt, order.Movie_id, order.Total_price, order.Full_name, order.Email, order.Phone_number, order.Payment_methode_id, order.Paid, order.Date, order.Time, order.Cinema_id}
	var orderId int
	var cinemaId int
	err := o.db.QueryRow(ctx, query, values...).Scan(&orderId, &cinemaId)
	if err != nil {
		return err
	}

	// melakukan looping berdasarkan inputan seats
	for _, seatStr := range order.SeatStr {

		// mengambil seat id pada table seats berdasarkan kode seat dan schedule id 
		querySelectIdSeat := "SELECT id FROM seats WHERE seat = $1 AND schedule_id = $2"
		var seatId int
		errSelectId := o.db.QueryRow(ctx, querySelectIdSeat, seatStr, cinemaId).Scan(&seatId)
		if errSelectId != nil {
			return errSelectId
		} 
		
		// menambahkan order_id dan seat_id pada tabel asosiasi order_seats 
		queryInsertOrderSeats := "INSERT INTO order_seats (order_id, seat_id) VALUES ($1, $2)"
		_, err := o.db.Exec(ctx, queryInsertOrderSeats, orderId, seatId)
		if err != nil {
			return err
		}

		// menandai seat yang sudah dibeli
		queryUpdateSeat := "UPDATE seats SET sold = true WHERE id = $1"
		_, errUpdate := o.db.Exec(ctx, queryUpdateSeat, seatId)

		if errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}

// repository get order history user
func (o *OrdersRepository) GetOrderHistory(ctx context.Context, IdInt int) ([]models.OrdersStruct, error) {

	// mengambil data dari beberapa tabel yang di joinkan (orders, payment_methode, movies, cinemas, order_seats, dan seats) berdasarkan user_id
	query := "SELECT o.user_id, o.total_price, o.paid, o.date, o.time, pm.name, m.title, c.image_path, array_agg(s.seat) FROM orders o JOIN payment_methode pm ON o.payment_methode_id = pm.id JOIN movies m ON m.id = o.movie_id JOIN cinemas c ON o.cinema_id = c.id JOIN order_seats os ON os.order_id = o.id join seats s ON s.id = os.seat_id WHERE o.user_id = $1 GROUP BY o.user_id, o.total_price, o.paid, o.date, o.time, pm.name, m.title, c.image_path"
	values := []any{IdInt}
	rows, err := o.db.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	} 
	
	defer rows.Close()

	var result []models.OrdersStruct
	for rows.Next() {
		var order models.OrdersStruct
		err := rows.Scan(&order.User_id, &order.Total_price, &order.Paid, &order.Date, &order.Time, &order.Payment_methode, &order.Title, &order.Cinema_path, &order.SeatStr)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}

	return result, nil
}
