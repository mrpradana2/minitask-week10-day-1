package repositories

import (
	"context"
	"tikcitz-app/internals/models"
	"tikcitz-app/pkg"

	"github.com/gin-gonic/gin"
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