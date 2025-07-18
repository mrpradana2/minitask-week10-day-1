package repositories

import (
	"context"
	"errors"
	"log"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersRepository struct {
	db *pgxpool.Pool
}

func NewOrdersRepository(db *pgxpool.Pool) *OrdersRepository {
	return &OrdersRepository{db: db}
}

// repository create order (fix)
func (o *OrdersRepository) CreateOrder(ctx context.Context, order models.OrdersStr, IdInt int) error {

	// awali dengan db transaction
	tx, err := o.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil
	}

	// rollback jika gagal menjalankan transaksi
	defer tx.Rollback(ctx)

	// insert to table orders and returning id and cinema_id
	query := "insert into orders(user_id, schedule_id, payment_methode_id, date, time, total_price, full_name, email, phone_number, paid) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id" //
	values := []any{IdInt, order.ScheduleId, order.PaymentMethodeId, order.Date, order.Time, order.TotalPrice, order.FullName, order.Email, order.PhoneNumber, order.Paid}
	var orderId int
	if err := tx.QueryRow(ctx, query, values...).Scan(&orderId); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// berikan poin kepada user yang membeli
	// update point user di table profiles
	queryAddPointUser := "UPDATE profiles SET point = point + 50 WHERE user_id = $1"
	if _, err := tx.Exec(ctx, queryAddPointUser, IdInt); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// build dinamic query untuk mengambil seat_id dari table seats 
	querySelectSeats, seats := utils.GetIdTable("seats", "kode", order.Seats)

	// mengeksekusi query select seat_id yang sudah di build
	rows, err := tx.Query(ctx, querySelectSeats, seats...)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}
	
	defer rows.Close()

	idSeats := []any{orderId}
	for rows.Next() {
		var idSeat int
		if err := rows.Scan(&idSeat); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return err
		}
		idSeats = append(idSeats, idSeat)
	}

	if len(order.Seats) != len(idSeats) - 1 {
		log.Println("[ERROR] ; ", errors.New("invalid choose seat"))
		return errors.New("kursi yang anda pesan tidak valid")
	}

	// menambahkan order_id dan order seat ke table asosiasi order_seats
	// melakukan build untuk query insert order_seats
	queryInsertOrderSeats := utils.InsertTableAssoc("order_seats", "order_id", "seat_id", idSeats)

	// mengeksekusi query insert order_seats yang sudah di build
	if _, err := tx.Exec(ctx, queryInsertOrderSeats, idSeats...); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// jangan lupa commit
	if err := tx.Commit(ctx); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	} 

	return nil
}

// repository get order history user (fix)
func (o *OrdersRepository) GetOrderHistory(ctx context.Context, IdInt int) ([]models.OrdersStr, error) {

	// mengambil data dari beberapa tabel yang di joinkan (orders, payment_methode, movies, cinemas, order_seats, dan seats) berdasarkan user_id
	query := "select o.id, m.title, c.image_path, o.date, o.time, o.total_price, o.paid, array_agg(s2.kode) from orders o join schedule s on o.schedule_id = s.id join payment_methode pm on pm.id = o.payment_methode_id join cinemas c on s.cinema_id = c.id join movies m on m.id = s.movie_id join users u on u.id = o.user_id join order_seats os on os.order_id = o.id join seats s2 on s2.id = os.seat_id where u.id = $1 group by o.id, c.image_path, m.title, o.date, o.time, o.total_price, o.paid order by o.create_at desc"
	values := []any{IdInt}
	rows, err := o.db.Query(ctx, query, values...)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	} 
	
	defer rows.Close()

	var result []models.OrdersStr
	for rows.Next() {
		var order models.OrdersStr
		err := rows.Scan(&order.Id, &order.Title, &order.ImagePath, &order.Date, &order.Time, &order.TotalPrice, &order.Paid, &order.Seats)
		if err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		result = append(result, order)
	}

	return result, nil
}

// repository get order by id (fix)
func (o *OrdersRepository) GetOrderById(ctx context.Context, userId, orderId int) ([]models.OrdersStr, error) {
	// mengambil data dari beberapa tabel yang di joinkan (orders, payment_methode, movies, cinemas, order_seats, dan seats) berdasarkan user_id
	query := "select o.id, m.title, c.image_path, o.date, o.time, o.total_price, o.paid, array_agg(s2.kode) from orders o join schedule s on o.schedule_id = s.id join payment_methode pm on pm.id = o.payment_methode_id join cinemas c on s.cinema_id = c.id join movies m on m.id = s.movie_id join users u on u.id = o.user_id join order_seats os on os.order_id = o.id join seats s2 on s2.id = os.seat_id where u.id = $1 and o.id = $2 group by o.id, c.image_path, m.title, o.date, o.time, o.total_price, o.paid order by o.create_at desc"
	values := []any{userId, orderId}
	rows, err := o.db.Query(ctx, query, values...)
	if err != nil {
		log.Println("[ERROR] ; ", err.Error())
		return nil, err
	} 
	
	defer rows.Close()

	var result []models.OrdersStr
	for rows.Next() {
		var order models.OrdersStr
		err := rows.Scan(&order.Id, &order.Title, &order.ImagePath, &order.Date, &order.Time, &order.TotalPrice, &order.Paid, &order.Seats)
		if err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		result = append(result, order)
	}

	return result, nil
}