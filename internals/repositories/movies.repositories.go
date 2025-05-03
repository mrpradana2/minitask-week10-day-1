package repositories

import (
	"context"
	"log"
	"tikcitz-app/internals/models"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MoviesRepository struct{
	db *pgxpool.Pool
}

func NewMoviesRepository(db *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{db: db}
}

// repository get movie all (fix)
func (u *MoviesRepository) GetMovies(ctx context.Context) ([]models.MoviesStruct, error) {

	// query dengan menggunakan CTE untuk mengambil all movie dan melakukan join dengan tabel movies_genres dan genres untuk mendapatkan genre list, serta dan hasilnya di joinkan dengan tabel movie_casts dan casts untuk mengambil cats list 
	query := `with table_movie_genres as (select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(g.name) as "genres" from movies m join movie_genres mg on m.id = mg.movie_id join genres g on g.id = mg.genre_id group by m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name) select t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres, array_agg(c.name) from table_movie_genres t join movie_casts mc on t.id = mc.movie_id join casts c on c.id = mc.cast_id group by t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres order by t.id asc`
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct

	// melakukan loop untuk memasukkan setiap movie ke variable result
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_movie, &movies.Duration, &movies.Director_name, &movies.Genres, &movies.Casts); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	
	return result, nil
}

// repository add movie (fix)
func (u *MoviesRepository) AddMovie(ctx context.Context, title, filePath, overview, directorName, location string, releaseDate, date time.Time, times []time.Time, duration, price int, genres, casts []string, cinemaIds []int) (error) {

	// menambahakan data movie baru dengan mereturn kan id movie yang baru dibuat
	query := "INSERT INTO movies (title, image_path, overview, release_date, director_name, duration) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	values := []any{title, filePath, overview, releaseDate, directorName, duration}
	var movieId int
	err := u.db.QueryRow(ctx, query, values...).Scan(&movieId)
	if err != nil {
		return err
	}

	// melakukan looping untuk mengisi data genres
	for _, genre := range genres {

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

	for _, cast :=range casts {
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
	for _, cinema := range cinemaIds {

		// lakukan looping untuk memasukkan jadwal berdasarkan time
		for _, time := range times {
			queryInsertSchedule := "INSERT INTO schedule (cinema_id, movie_id, location, date, time, price) VALUES ($1, $2, $3, $4, $5, $6)"

			if _, err := u.db.Exec(ctx, queryInsertSchedule, cinema, movieId, location, date, time, price); err != nil {
				return err
			}

		}
	}
	return nil
}

// repository update movie (fix)
func (u *MoviesRepository) UpdateMovie(ctx context.Context, title, filePath, overview, directorName string, releaseDate time.Time, duration int, genres, casts []string, idInt int) (pgconn.CommandTag, error) {

	// melakukan update movie berdasarkan id movie
	query := "update movies set title = $1, image_path = $2, overview = $3, release_date = $4, director_name = $5, duration = $6, modified_at = $7 where id = $8;"

	values := []any{title, filePath, overview, releaseDate, directorName, duration, time.Now(), idInt}

	cmd, err := u.db.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	// memperbaharui movie_genres
	// melakukan delete pada tabel asosiasi movie_genre berdasarkan movie id
	queryDelMovieGenre := "DELETE FROM movie_genres WHERE movie_id = $1"
	_, errDelMovieGenre := u.db.Exec(ctx, queryDelMovieGenre, idInt)
	if errDelMovieGenre != nil {
		return pgconn.CommandTag{}, nil
	}
	log.Println("GENRESSS : ", genres)
	for _, genre := range genres {
		log.Println("[error bro]")
		// menambahkan genre baru jika belum terdaftar
		queryGenres := "INSERT INTO genres (name) VALUES ($1) ON CONFLICT (name) DO NOTHING"
		_, err := u.db.Exec(ctx, queryGenres, genre)
		if err != nil {
			return pgconn.CommandTag{}, err
		}

		// ambil genre id
		var genreId int
		queryGenreId := "SELECT id FROM genres WHERE name = $1"
        err = u.db.QueryRow(ctx, queryGenreId, genre).Scan(&genreId)
        if err != nil {
            return pgconn.CommandTag{}, err
        }

		// tambahkan movie id dan genre id ke tabel asosiasi movie_genre
		queryMovieGenre := "INSERT INTO movie_genres (movie_id, genre_id) VALUES ($1, $2)"
		_, err = u.db.Exec(ctx, queryMovieGenre, idInt, genreId)
        if err != nil {
            return pgconn.CommandTag{}, err
        }
	}

		// memperbarui movie_casts
		// melakukan delete pada tabel asosiasi movie_casts berdasarkan movie id
		queryDelMovieCasts := "DELETE FROM movie_casts WHERE movie_id = $1"
		_, errDelMovieCasts := u.db.Exec(ctx, queryDelMovieCasts, idInt)
		if errDelMovieCasts != nil {
			return pgconn.CommandTag{}, nil
		}

	 	for _, cast := range casts {
		// menambahkan genre baru jika belum terdaftar
		queryGenres := "INSERT INTO casts (name) VALUES ($1) ON CONFLICT (name) DO NOTHING"
		_, err := u.db.Exec(ctx, queryGenres, cast)
		if err != nil {
			return pgconn.CommandTag{}, err
		}

		// ambil genre id
		var castId int
		queryGenreId := "SELECT id FROM casts WHERE name = $1"
        err = u.db.QueryRow(ctx, queryGenreId, cast).Scan(&castId)
        if err != nil {
            return pgconn.CommandTag{}, err
        }

		// tambahkan movie id dan genre id ke tabel asosiasi movie_genre
		queryMovieGenre := "INSERT INTO movie_casts (movie_id, cast_id) VALUES ($1, $2)"
		_, err = u.db.Exec(ctx, queryMovieGenre, idInt, castId)
        if err != nil {
            return pgconn.CommandTag{}, err
        }
	}

	return cmd, nil
}

// repository delete movie (fix)
func (u *MoviesRepository) DeleteMovie(ctx context.Context, idInt int) (pgconn.CommandTag, error) { 

	// melakukan delele movie_genres berdasarkan movie_id
	queryMovieGenres := "DELETE FROM movie_genres WHERE movie_id = $1"
	_, errMovieGenres := u.db.Exec(ctx, queryMovieGenres, idInt)
	if errMovieGenres != nil {
		return pgconn.CommandTag{}, nil
	}

	// melakukan delete movie_casts berdasarkan movie_id
	queryMovieCasts := "DELETE FROM movie_casts WHERE movie_id = $1"
	_, errMovieCasts := u.db.Exec(ctx, queryMovieCasts, idInt)
	if errMovieCasts != nil {
		return pgconn.CommandTag{}, nil
	}

	// melakukan delete schedule berdasarkan movie_id
	querySchedule := "DELETE FROM schedule WHERE movie_id = $1"
	_, errSchedule := u.db.Exec(ctx, querySchedule, idInt)
	if errSchedule != nil {
		return pgconn.CommandTag{}, nil
	}

	// melakukan delete movie berdasarkan movie id
	queryMovie := "DELETE FROM movies WHERE id = $1"
	cmd, err := u.db.Exec(ctx, queryMovie, idInt)
	if err != nil {
		return pgconn.CommandTag{}, nil
	}

	return cmd, nil
}

// repository get upcoming movie (fix)
func (u *MoviesRepository) GetMovieUpcoming(ctx context.Context) ([]models.MoviesStruct, error) {
	
	// query dengan menggunakan CTE untuk mengambil all movie dan melakukan join dengan tabel movies_genres dan genres untuk mendapatkan genre list, serta dan hasilnya di joinkan dengan tabel movie_casts dan casts untuk mengambil cats list 
	query := `with table_movie_genres as (select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(g.name) as "genres" from movies m join movie_genres mg on m.id = mg.movie_id join genres g on g.id = mg.genre_id where m.release_date > now() group by m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name) select t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres, array_agg(c.name) from table_movie_genres t join movie_casts mc on t.id = mc.movie_id join casts c on c.id = mc.cast_id group by t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres order by t.id asc`
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct

	// melakukan loop untuk memasukkan setiap movie ke variable result
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_movie, &movies.Duration, &movies.Director_name, &movies.Genres, &movies.Casts); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	
	return result, nil

	// // mengambil data movie 
	// query := "SELECT m.id, m.title, sm.status, m.release_date, m.overview, m.image_path, m.duration, m.director_name, m.casts, ARRAY_AGG(g.genre_name) FROM movies m JOIN status_movie sm ON m.status_movie_id = sm.id JOIN movie_genre mg ON mg.movie_id = m.id JOIN genres g ON g.id = mg.genre_id WHERE status_movie_id = 1 GROUP BY m.id, sm.status"

	// rows, err := u.db.Query(ctx, query)
	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()
	// var result []models.MoviesStruct
	// for rows.Next() {
	// 	var movies models.MoviesStruct
	// 	if err := rows.Scan(&movies.Id, &movies.Title, &movies.Status_movie, &movies.Release_date, &movies.Overview, &movies.Image_path, &movies.Duration, &movies.Director_name, &movies.Casts, &movies.Genres); err != nil {
	// 		return nil, err
	// 	}
	// 	result = append(result, movies)
	// }
	// return result, nil
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

// repository get detail movie (fix)
func (u *MoviesRepository) GetDetailMovie(ctx context.Context, movies models.MoviesStruct, IdInt int) ([]models.MoviesStruct, error) {

	// query dengan menggunakan CTE untuk mengambil detail movie detail dengan id tertentu dan melakukan join dengan tabel movies_genres dan genres untuk mendapatkan genre list, serta dan hasilnya di joinkan dengan tabel movie_casts dan casts untuk mengambil cats list 
	query := `with table_movie_genres as (select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(g.name) as "genres" from movies m join movie_genres mg on m.id = mg.movie_id join genres g on g.id = mg.genre_id where m.id = $1 group by m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name) select t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres, array_agg(c.name) from table_movie_genres t join movie_casts mc on t.id = mc.movie_id join casts c on c.id = mc.cast_id group by t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres`
	// values := []any{IdInt}
	rows, err := u.db.Query(ctx, query, IdInt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct

	// melakukan loop untuk memasukkan setiap movie ke variable result
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_movie, &movies.Duration, &movies.Director_name, &movies.Genres, &movies.Casts); err != nil {
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