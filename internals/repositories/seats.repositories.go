package repositories

import (
	"context"
	"log"
	"tikcitz-app/internals/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatsRepository struct {
	db *pgxpool.Pool
}

func NewSeatsRepository(db *pgxpool.Pool) *SeatsRepository {
	return &SeatsRepository{db: db}
}

// repository get available seats
func (s *SeatsRepository) GetSeatsAvailable(ctx context.Context, seat models.SeatsStruct, id int) (models.ResultSeat, error) {
	
	tx, err := s.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.ResultSeat{}, err
	}

	defer tx.Rollback(ctx)

	// mengambil seat available
	// query mengambil table seats, dengan ketentuan tidak mengambil hasil seat yang terdapat pada sub query
	// didalam sub query mengambil table orders yang di joinkan dengan table asosiasi order_seats dimana diambil berdasarkan schedule id tertentu dan jika didalam table asosisasi terdapat id yang sama dengan seat id yang berarti kursi tersebut sudah pernah dipesan 
	query := "select s.id, s.kode from seats s where not exists (select s.kode from orders o2 join order_seats os2 on o2.id = os2.order_id where o2.schedule_id = $1 and os2.seat_id = s.id)"

	// mengambil data hasil query
	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.ResultSeat{}, err
	}

	defer rows.Close()
	var result []models.SeatsStruct

	// melakukan loop untuk memasukkan setiap movie ke variable result
	for rows.Next() {
		var seat models.SeatsStruct
		if err := rows.Scan(&seat.Id, &seat.Seat); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return models.ResultSeat{}, err
		}
		result = append(result, seat)
	}

	// mengambil data schedule id, title, cinema, tanggal dan waktu berdasarkan schedule id
	queryGetSchedule := "select s.id, m.title, c.name, s.date, s.time from schedule s join movies m on s.movie_id = m.id join cinemas c on c.id = s.cinema_id where s.id = $1"

	// menyiapkan variable untuk menampung data schedule dan seat available
	var availableSeat models.ResultSeat

	// menjalankan query row untuk mengakses database untuk mendapatkan data schedule
	if err := tx.QueryRow(ctx, queryGetSchedule, id).Scan(&availableSeat.Id, &availableSeat.Title, &availableSeat.Cinema, &availableSeat.Date, &availableSeat.Time); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.ResultSeat{}, err
	}

	// memasukkan hasil query get available seat kedalam variable availableSeat 
	for _, res := range result {
		availableSeat.Seats = append(availableSeat.Seats, res.Seat)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("[ERROR] ; ", err.Error())
		return models.ResultSeat{}, err
	}

	return availableSeat, nil
}