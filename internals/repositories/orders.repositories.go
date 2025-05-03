package repositories

import (
	"context"
	"errors"
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
func (o *OrdersRepository) CreateOrder(ctx context.Context, order models.OrdersStr, IdInt int) error {

	// insert to table orders and returning id and cinema_id
	query := "insert into orders(user_id, schedule_id, payment_methode_id, date, time, total_price, full_name, email, phone_number, paid) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id" //
	values := []any{IdInt, order.ScheduleId, order.PaymentMethodeId, order.Date, order.Time, order.TotalPrice, order.FullName, order.Email, order.PhoneNumber, order.Paid}
	var orderId int
	err := o.db.QueryRow(ctx, query, values...).Scan(&orderId)
	if err != nil {
		return err
	}

	// berikan poin kepada user yang membeli
	// mengambil nilai point user dari tabel profile
	querySelectPoint := "SELECT point FROM profiles WHERE user_id = $1"
	var point int
	if err := o.db.QueryRow(ctx, querySelectPoint, IdInt).Scan(&point); err != nil {
		return err
	}

	// menambahkan point +50 untuk setiap kali order
	point = point + 50

	// melakukan update point terbaru
	queryUpdatePoint := "UPDATE profiles SET point = $1 WHERE user_id = $2"
	if	_, err := o.db.Exec(ctx, queryUpdatePoint, point, IdInt); err != nil {
		return err
	}

	// melakukan looping berdasarkan inputan seats
	for _, seatStr := range order.Seats {

		// mengambil seat id pada table seats berdasarkan kode seat dan schedule id 
		querySelectIdSeat := "SELECT id FROM seats WHERE kode = $1"
		var seatId int

		if err := o.db.QueryRow(ctx, querySelectIdSeat, seatStr).Scan(&seatId); err != nil {
			return err
		} 

		if seatId == 0 {
			return errors.New("kursi tidak tersedia")
		}

		// menambahkan order_id dan seat_id pada tabel asosiasi order_seats 
		queryInsertOrderSeats := "INSERT INTO order_seats (order_id, seat_id) VALUES ($1, $2)"
		if _, err := o.db.Exec(ctx, queryInsertOrderSeats, orderId, seatId); err != nil {
			return err
		}
	}
	return nil
}

// repository get order history user
func (o *OrdersRepository) GetOrderHistory(ctx context.Context, IdInt int) ([]models.OrdersStr, error) {

	// mengambil data dari beberapa tabel yang di joinkan (orders, payment_methode, movies, cinemas, order_seats, dan seats) berdasarkan user_id
	query := "select o.id, m.title, c.image_path, o.date, o.time, o.total_price, o.paid, array_agg(s2.kode) from orders o join schedule s on o.schedule_id = s.id join payment_methode pm on pm.id = o.payment_methode_id join cinemas c on s.cinema_id = c.id join movies m on m.id = s.movie_id join users u on u.id = o.user_id join order_seats os on os.order_id = o.id join seats s2 on s2.id = os.seat_id where u.id = $1 group by o.id, c.image_path, m.title, o.date, o.time, o.total_price, o.paid"
	values := []any{IdInt}
	rows, err := o.db.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	} 
	
	defer rows.Close()

	var result []models.OrdersStr
	for rows.Next() {
		var order models.OrdersStr
		err := rows.Scan(&order.Id, &order.Title, &order.ImagePath, &order.Date, &order.Time, &order.TotalPrice, &order.Paid, &order.Seats)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}

	return result, nil
}
