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
func (s *SeatsRepository) GetSeatsAvailable(ctx context.Context, seat models.SeatsStruct, id int) ([]models.SeatsStruct, error) {
	
	// mengambil seat available
	query := "select s.id, s.kode from seats s where not exists (select s.kode from orders o2 join order_seats os2 on o2.id = os2.order_id where o2.schedule_id = $1 and os2.seat_id = s.id)"

	rows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.SeatsStruct

	// melakukan loop untuk memasukkan setiap movie ke variable result
	for rows.Next() {
		var seat models.SeatsStruct
		if err := rows.Scan(&seat.Id, &seat.Seat); err != nil {
			return nil, err
		}
		result = append(result, seat)
	}

	return result, nil
}