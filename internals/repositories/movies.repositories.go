package repositories

import (
	"context"
	"tikcitz-app/pkg"
	"time"

	"github.com/gin-gonic/gin"
)

type moviesStruct struct {
	Id              int       `json:"id" form:"id,omitempty"`
	Title           string    `json:"title" form:"title"`
	Image_path      string    `json:"image_path" form:"image_path"`
	Overview        string    `json:"overview" form:"overview"`
	Release_date    time.Time `json:"release_date" form:"release_date"`
	Director_name   string    `json:"director_name" form:"director_name"`
	Duration        int       `json:"duration" form:"duration"`
	Casts           []string  `json:"casts" form:"casts"`
	Status_movie_id int       `json:"status_movie_id" form:"status_movie_id"`
	Genres          []string  `json:"genres" form:"genres"`
	Status_movie    string    `json:"status_movie" form:"status_movie"`
}

type MoviesRepository struct{}

var MovieRepo *MoviesRepository

func NewMoviesRepository() {
	MovieRepo = &MoviesRepository{}
}

func (u *MoviesRepository) GetMovies(ctx *gin.Context) ([]moviesStruct, error) {
	query := "SELECT m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, array_agg(g.genre_name) from movies m join movie_genre mg on m.id = mg.movie_id join genres g on mg.genre_id = g.id group by m.id;"
	rows, err := pkg.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []moviesStruct
	for rows.Next() {
		var movies moviesStruct
		if err := rows.Scan(&movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	
	return result, nil
}