package repositories

import (
	"context"
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
func (s *SeatsRepository) GetSeatsAvailable(ctx context.Context, cinema models.MoviesStruct, query string) ([]models.SeatsStruct, error) {
	
	// mengambil data cinema di tabel cinema
	queryGetCinema := "SELECT id, cinema_name, image_path FROM cinemas"
	rows, errGetcinema := s.db.Query(ctx, queryGetCinema)
	if errGetcinema != nil {
		return []models.SeatsStruct{}, errGetcinema
	}

	defer rows.Close()
	var findCinema []models.CinemaStruct
	for rows.Next() {
		var cinema models.CinemaStruct
		err := rows.Scan(&cinema.Id, &cinema.Cinema_name, &cinema.Image_path)
		if err != nil {
			return []models.SeatsStruct{}, err
		}

		// mengecek jika cinema_name(dari query) ada dengan data cinema_name di database maka masukkan ke variable findcinema 
		if cinema.Cinema_name == query {
			findCinema = append(findCinema, cinema)
		}
	}

	if len(findCinema) == 0 {
		return []models.SeatsStruct{}, nil
	}

	// mengambil data schedule berdasarkan cinema id dan kursi yang sudah terjual
	queryGetSeats := "SELECT s.id, s.seat, s.sold, c.cinema_name FROM seats s JOIN schedule s2 ON s.schedule_id = s2.id JOIN cinemas c ON c.id = s2.cinema_id WHERE s2.cinema_id = $1 AND s.sold = false"
	values := []any{findCinema[0].Id}
	rowsSeat, errGetSeat := s.db.Query(ctx, queryGetSeats, values...)
	if errGetSeat != nil {
		return []models.SeatsStruct{}, errGetSeat
	}
	
	defer rowsSeat.Close()
	var seats []models.SeatsStruct
	
	for rowsSeat.Next() {
		var seat models.SeatsStruct
		err := rowsSeat.Scan(&seat.Id, &seat.Seat, &seat.Sold, &seat.Cinema)
		if err != nil {
			return []models.SeatsStruct{}, err
		}
		seats = append(seats, seat)
	}
	
	return seats, nil
}