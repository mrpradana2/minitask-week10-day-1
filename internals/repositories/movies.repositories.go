package repositories

import (
	"context"
	"tikcitz-app/internals/models"
	"tikcitz-app/pkg"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type MoviesRepository struct{}

var MovieRepo *MoviesRepository

func NewMoviesRepository() {
	MovieRepo = &MoviesRepository{}
}

func (u *MoviesRepository) GetMovies(ctx *gin.Context) ([]models.MoviesStruct, error) {
	query := "SELECT m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, array_agg(g.genre_name) from movies m join movie_genre mg on m.id = mg.movie_id join genres g on mg.genre_id = g.id group by m.id;"
	rows, err := pkg.DB.Query(context.Background(), query)
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

func (u *MoviesRepository) AddMovie(ctx *gin.Context, newDataMovie models.MoviesStruct) (pgconn.CommandTag, error) {
	query := "INSERT INTO movies (title, image_path, overview, release_date, director_name, duration, casts, status_movie_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	values := []any{newDataMovie.Title, newDataMovie.Image_path, newDataMovie.Overview, newDataMovie.Release_date, newDataMovie.Director_name, newDataMovie.Duration, newDataMovie.Casts, newDataMovie.Status_movie_id}
	cmd, err :=pkg.DB.Exec(ctx.Request.Context(), query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}