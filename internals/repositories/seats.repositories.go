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
	
	// mengambil seat available
	query := "select s.id, s.kode from seats s where not exists (select s.kode from orders o2 join order_seats os2 on o2.id = os2.order_id where o2.schedule_id = $1 and os2.seat_id = s.id)"

	rows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return models.ResultSeat{}, err
	}

	defer rows.Close()
	var result []models.SeatsStruct

	// melakukan loop untuk memasukkan setiap movie ke variable result
	for rows.Next() {
		var seat models.SeatsStruct
		if err := rows.Scan(&seat.Id, &seat.Seat); err != nil {
			return models.ResultSeat{}, err
		}
		result = append(result, seat)
	}

	queryGetSchedule := "select s.id, m.title, c.name, s.date, s.time from schedule s join movies m on s.movie_id = m.id join cinemas c on c.id = s.cinema_id where s.id = $1"
	var schedule models.ResultSeat
	if err := s.db.QueryRow(ctx, queryGetSchedule, id).Scan(&schedule.Id, &schedule.Title, &schedule.Cinema, &schedule.Date, &schedule.Time); err != nil {
		return models.ResultSeat{}, err
	}

	// var hasil models.ResultSeat

	for _, res := range result {
		schedule.Seats = append(schedule.Seats, res.Seat)
	}

	log.Println(schedule)

	return schedule, nil
}