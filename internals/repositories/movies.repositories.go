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

	// mengambil data movie yang di join dengan table genre
	query := "SELECT m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, array_agg(g.genre_name) from movies m join movie_genre mg on m.id = mg.movie_id join genres g on mg.genre_id = g.id group by m.id"
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	
	return result, nil
}

// repository add movie
func (u *MoviesRepository) AddMovie(ctx context.Context, newDataMovie models.MoviesStruct) (error) {

	// menambahakan data movie baru dengan mereturn kan id movie yang baru dibuat
	query := "INSERT INTO movies (title, image_path, overview, release_date, director_name, duration) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	values := []any{newDataMovie.Title, newDataMovie.Image_path, newDataMovie.Overview, newDataMovie.Release_date, newDataMovie.Director_name, newDataMovie.Duration}
	var movieId int
	err := u.db.QueryRow(ctx, query, values...).Scan(&movieId)
	if err != nil {
		return err
	}

	// melakukan looping untuk mengisi data genres
	for _, genre := range newDataMovie.Genres {

		// menambahkan genre baru jika belum terdaftar
		queryGenres := "INSERT INTO genres (name) VALUES ($1) ON CONFLICT (name) DO NOTHING"
		_, err := u.db.Exec(ctx, queryGenres, genre)
		if err != nil {
			return err
		}

		// ambil genre id
		var genreId int
		queryGenreId := "SELECT id FROM genres WHERE name = $1"
        err = u.db.QueryRow(ctx, queryGenreId, genre).Scan(&genreId)
        if err != nil {
            return err
        }

		// tambahkan movie id dan genre id ke tabel asosiasi movie_genre
		queryMovieGenre := "INSERT INTO movie_genres (movie_id, genre_id) VALUES ($1, $2)"
		_, err = u.db.Exec(ctx, queryMovieGenre, movieId, genreId)
        if err != nil {
            return err
        }
	}

	for _, cast :=range newDataMovie.Casts {
		// menambahkan cast baru jika belum ada
		queryCast := "INSERT INTO casts(name) VALUES($1) ON CONFLICT (name) DO NOTHING"
		if _, err := u.db.Exec(ctx, queryCast, cast); err != nil {
			return err
		}

		// ambil cast id
		queryGetIdCast := "SELECT id FROM casts WHERE name = $1"
		var castId int
		if err = u.db.QueryRow(ctx, queryGetIdCast, cast).Scan(&castId); err != nil {
            return err
        }

		// tambahkan cast id dan movie id ke tabel asosiasi movie_casts
		queryMovieGenre := "INSERT INTO movie_casts (movie_id, cast_id) VALUES ($1, $2)"
		if _, err = u.db.Exec(ctx, queryMovieGenre, movieId, castId); err != nil {
            return err
        }
	}

	// tambahkan jadwal untuk movie ini
	// lakukan looping untuk memasukkan jadwal berdasarkan movie yang akan ditampilkan cinema 
	for _, cinema := range newDataMovie.Cinema_ids {

		// lakukan looping untuk memasukkan jadwal berdasarkan time
		for _, time := range newDataMovie.Times {
			queryInsertSchedule := "INSERT INTO schedule (cinema_id, movie_id, location, date, time, price) VALUES ($1, $2, $3, $4, $5, $6)"

			if _, err := u.db.Exec(ctx, queryInsertSchedule, cinema, movieId, newDataMovie.Location, newDataMovie.Date, time, newDataMovie.Price); err != nil {
				return err
			}

		}
	}
	return nil
}

// repository update movie
func (u *MoviesRepository) UpdateMovie(ctx context.Context, updateMovie *models.MoviesStruct, idInt int) (pgconn.CommandTag, error) {

	// melakukan update movie berdasarkan id movie
	query := "UPDATE movies SET title = $1, image_path = $2, overview = $3, release_date = $4, director_name = $5, duration = $6, casts = $7, status_movie_id = $8 WHERE id = $9"

	values := []any{updateMovie.Title, updateMovie.Image_path, updateMovie.Overview, updateMovie.Release_date, updateMovie.Director_name, updateMovie.Duration, updateMovie.Casts, updateMovie.Status_movie_id, idInt}

	cmd, err := u.db.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	// melakukan delete pada tabel asosiasi movie_genre berdasarkan movie id
	queryDelMovieGenre := "DELETE FROM movie_genre WHERE movie_id = $1"
	_, errDelMovieGenre := u.db.Exec(ctx, queryDelMovieGenre, idInt)
	if errDelMovieGenre != nil {
		return pgconn.CommandTag{}, nil
	}

	for _, genre := range updateMovie.Genres {
		// menambahkan genre baru jika belum terdaftar
		queryGenres := "INSERT INTO genres (genre_name) VALUES ($1) ON CONFLICT (genre_name) DO NOTHING"
		_, err := u.db.Exec(ctx, queryGenres, genre)
		if err != nil {
			return pgconn.CommandTag{}, err
		}

		// ambil genre id
		var genreId int
		queryGenreId := "SELECT id FROM genres WHERE genre_name = $1"
        err = u.db.QueryRow(ctx, queryGenreId, genre).Scan(&genreId)
        if err != nil {
            return pgconn.CommandTag{}, err
        }

		// tambahkan movie id dan genre id ke tabel asosiasi movie_genre
		queryMovieGenre := "INSERT INTO movie_genre (movie_id, genre_id) VALUES ($1, $2)"
		_, err = u.db.Exec(ctx, queryMovieGenre, idInt, genreId)
        if err != nil {
            return pgconn.CommandTag{}, err
        }
	}

	return cmd, nil
}

// repository delete movie
func (u *MoviesRepository) DeleteMovie(ctx context.Context, idInt int) (pgconn.CommandTag, error) {

	// melakukan delete movie berdasarkan movie id
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

	// mengambil data movie 
	query := "SELECT m.id, m.title, sm.status, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, ARRAY_AGG(g.genre_name) FROM movies m JOIN status_movie sm ON m.status_movie_id = sm.id JOIN movie_genre mg ON mg.movie_id = m.id JOIN genres g ON g.id = mg.genre_id WHERE status_movie_id = 1 GROUP BY m.id, sm.status"

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Status_movie, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	return result, nil
}

// repository get popular movie
func (u *MoviesRepository) GetMoviePopular(ctx context.Context) ([]models.MoviesStruct, error) {
	query := "SELECT m.id, m.title, sm.status, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, ARRAY_AGG(g.genre_name) FROM movies m JOIN status_movie sm ON m.status_movie_id = sm.id JOIN movie_genre mg ON mg.movie_id = m.id JOIN genres g ON g.id = mg.genre_id WHERE status_movie_id = 2 GROUP BY m.id, sm.status"

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Status_movie, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	return result, nil
}

// repository get detail movie
func (u *MoviesRepository) GetDetailMovie(ctx context.Context, movies models.MoviesStruct, IdInt int) ([]models.MoviesStruct, error) {

	// mengambil data movie dan melakukan join dengan tabel asosiasi movie_genre dan tabel genre untuk mengambil genre yang digabung menjadi array berdasarkan id movie 
	query := "SELECT m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, array_agg(g.genre_name) FROM movies m JOIN movie_genre mg ON m.id = mg.movie_id JOIN genres g ON mg.genre_id = g.id WHERE m.id = $1 GROUP BY m.id"
	values := []any{IdInt}
	rows, err := u.db.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}

	return result, nil
}

// repository get movie with pagination
func (u *MoviesRepository) GetMoviesWithPagination(ctx context.Context, movie models.MoviesStruct, offset int) ([]models.MoviesStruct, error) {

	// mengambil data movies menggunakan paginasi yang di join dengan tabel asosiasi movie_genre dan tabel genres untuk mengambil genre yang gabung menjadi array
	query := "SELECT m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, ARRAY_AGG(g.genre_name) FROM movies m JOIN movie_genre mg ON m.id = mg.movie_id JOIN genres g ON mg.genre_id = g.id GROUP BY m.id ORDER BY m.id ASC LIMIT 5 OFFSET $1;"
	values := []any{offset}
	rows, err := u.db.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	
	return result, nil
}