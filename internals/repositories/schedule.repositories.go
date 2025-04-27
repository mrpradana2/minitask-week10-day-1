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
func (s *ScheduleRepository) GetScheduleMovie(ctx context.Context, schedule *models.ScheduleStruct) ([]models.ScheduleStruct, error) {
	query := "SELECT s.id, c.cinema_name, s.date, s.times, s.location, s.price FROM schedule s JOIN cinemas c ON s.cinema_id = c.id"
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.ScheduleStruct
	for rows.Next() {
		var schedule models.ScheduleStruct
		err := rows.Scan(&schedule.Id, &schedule.Cinema, &schedule.Date, &schedule.Time, &schedule.Location, &schedule.Price)
		if err != nil {
			return nil, err
		}

		schedule.DateStr = schedule.Date.Format("2006-01-02")
		// schedule.TimeStr = schedule.Time.Format("15:04")

		result = append(result, schedule)
	}
	return result, nil
}
