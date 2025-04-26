package repositories

import (
	"context"
	"tikcitz-app/internals/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MoviesRepository struct{
	db *pgxpool.Pool
}

func NewMoviesRepository(db *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{db: db}
}

// repository get movie all
func (u *MoviesRepository) GetMovies(ctx context.Context) ([]models.MoviesStruct, error) {
	query := "SELECT m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, array_agg(g.genre_name) from movies m join movie_genre mg on m.id = mg.movie_id join genres g on mg.genre_id = g.id group by m.id;"
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	
	return result, nil
}

// repository add movie
func (u *MoviesRepository) AddMovie(ctx context.Context, newDataMovie models.MoviesStruct) (pgconn.CommandTag, error) {
	query := "INSERT INTO movies (title, image_path, overview, release_date, director_name, duration, casts, status_movie_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	values := []any{newDataMovie.Title, newDataMovie.Image_path, newDataMovie.Overview, newDataMovie.Release_date, newDataMovie.Director_name, newDataMovie.Duration, newDataMovie.Casts, newDataMovie.Status_movie_id}
	cmd, err := u.db.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}

// repository update movie
func (u *MoviesRepository) UpdateMovie(ctx context.Context, updateMovie *models.MoviesStruct, idInt int) (pgconn.CommandTag, error) {
	query := "UPDATE movies SET title = $1, image_path = $2, overview = $3, release_date = $4, director_name = $5, duration = $6, casts = $7, status_movie_id = $8 WHERE id = $9"

	values := []any{updateMovie.Title, updateMovie.Image_path, updateMovie.Overview, updateMovie.Release_date, updateMovie.Director_name, updateMovie.Duration, updateMovie.Casts, updateMovie.Status_movie_id, idInt}

	cmd, err := u.db.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}

// repository delete movie
func (u *MoviesRepository) DeleteMovie(ctx context.Context, idInt int) (pgconn.CommandTag, error) {
	query := "DELETE FROM movies WHERE id = $1"
	values := []any{idInt}
	cmd, err := u.db.Exec(ctx, query, values...)

	if err != nil {
		return pgconn.CommandTag{}, nil
	}

	return cmd, nil
}

// repository get upcoming movie
func (u *MoviesRepository) GetMovieUpcoming(ctx context.Context) ([]models.MoviesStruct, error) {
	query := "SELECT m.title, sm.status, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts FROM movies m JOIN status_movie sm ON m.status_movie_id = sm.id WHERE status_movie_id = 1"

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Title, &movies.Status_movie, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	return result, nil
}

// repository get popular movie
func (u *MoviesRepository) GetMoviePopular(ctx context.Context) ([]models.MoviesStruct, error) {
	query := "SELECT m.title, sm.status, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts FROM movies m JOIN status_movie sm ON m.status_movie_id = sm.id WHERE status_movie_id = 2"

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Title, &movies.Status_movie, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	return result, nil
}