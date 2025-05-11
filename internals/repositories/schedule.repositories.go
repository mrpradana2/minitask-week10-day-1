package repositories

import (
	"context"
	"tikcitz-app/internals/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// repository get schedule
func (s *ScheduleRepository) GetScheduleMovie(ctx context.Context, schedule *models.ScheduleStruct, idInt int) ([]models.ScheduleStruct, error) {

	// mengambil data schedule dengan join tabel cinema untuk mendapatkan informasi cinema  
	// query := "select m.id, m.title, c.name, c.image_path, s.date, array_agg(s.time), s.location, s.price, array_agg(s.id), from schedule s join cinemas c on s.cinema_id = c.id join movies m on s.movie_id = m.id where m.id = $1 group by m.id, m.title, c.name, c.image_path, s.date, s.location, s.price"
	query := "select s.id, m.id, m.title, c.name, c.image_path, s.date, s.time, s.location, s.price from schedule s join cinemas c on s.cinema_id = c.id join movies m on s.movie_id = m.id where m.id = $1"
	rows, err := s.db.Query(ctx, query, idInt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.ScheduleStruct
	for rows.Next() {
		var schedule models.ScheduleStruct
		err := rows.Scan(&schedule.Id, &schedule.MovieId, &schedule.Title, &schedule.Cinema, &schedule.CinemaPathImage, &schedule.Date, &schedule.Time, &schedule.Location, &schedule.Price)
		if err != nil {
			return nil, err
		}

		schedule.DateStr = schedule.Date.Format("2006-01-02")
		// schedule.TimeStr = schedule.Time.Format("14.30")

		result = append(result, schedule)
	}
	return result, nil
}
