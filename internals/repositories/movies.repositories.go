package repositories

import (
	"context"
	"log"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/utils"
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

	// mengambil data movies dengan menjoinkan table movies dengan tabel asosiasi movie_genres dan genres untuk mendapatkan genre movie, dan menjoinkan tabel asosiasi movie_casts dan casts untuk mendapatkan data cast, kolom genres dan casts digabungkan dengan array aggregat dan agar tidak duplikat tambahkan distinct 
	query := `select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(distinct g.name) as "genres", array_agg(distinct c.name) from movies m join movie_genres mg on mg.movie_id = m.id join genres g on mg.genre_id = g.id join movie_casts mc on mc.movie_id = m.id join casts c on c.id = mc.cast_id group by m.id order by m.id asc`
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

	tx, errBegin := u.db.Begin(ctx)
	if errBegin != nil {
		return errBegin
	}

	defer tx.Rollback(ctx)

	// menambahakan data movie baru dengan mereturn kan id movie yang baru dibuat
	query := "INSERT INTO movies (title, image_path, overview, release_date, director_name, duration) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	values := []any{title, filePath, overview, releaseDate, directorName, duration}
	var movieId int
	err := tx.QueryRow(ctx, query, values...).Scan(&movieId)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// memasukkan data genre yang belum terdaftar di table genres
	var genresAny []any

	for _, genre := range genres {
		genresAny = append(genresAny, genre)
	}

	queryGenres := utils.AddList("genres", "name", genres)
	if _, err := tx.Exec(ctx, queryGenres, genresAny...); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// build dinamic query untuk mengambil genre_id dari table genres 
	querySelectGenres, genresId := utils.GetIdTable("genres", "name", genres)

	// mengeksekusi query select seat_id yang sudah di build
	rows, err := tx.Query(ctx, querySelectGenres, genresId...)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	
	defer rows.Close()

	idGenres := []any{movieId}
	for rows.Next() {
		var idGenre int
		if err := rows.Scan(&idGenre); err != nil {
			return err
		}
		idGenres = append(idGenres, idGenre)
	}

	// menambahkan movie_id dan genre seat ke table asosiasi movie_genres
	// melakukan build untuk query insert movie_genres
	queryInsertMovieGenres := utils.InsertTableAssoc("movie_genres", "movie_id", "genre_id", idGenres)
	// mengeksekusi query insert movie_genres yang sudah di build
	if _, err := tx.Exec(ctx, queryInsertMovieGenres, idGenres...); err != nil {
		log.Println("[ERROR] :", err.Error())
		return err
	}

	// memasukkan data cast yang belum terdaftar di table genres
	var castsAny []any

	for _, cast := range casts {
		castsAny = append(castsAny, cast)
	}

	queryCasts := utils.AddList("casts", "name", casts)
	if _, err := tx.Exec(ctx, queryCasts, castsAny...); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// build dinamic query untuk mengambil cast_id dari table casts 
	querySelectCasts, castsId := utils.GetIdTable("casts", "name", casts)

	// mengeksekusi query select cast_id yang sudah di build
	rowsCasts, err := tx.Query(ctx, querySelectCasts, castsId...)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		return err
	}
	
	defer rowsCasts.Close()

	idCasts := []any{movieId}
	for rowsCasts.Next() {
		var idCast int
		if err := rowsCasts.Scan(&idCast); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return err
		}
		idCasts = append(idCasts, idCast)
	}

	// menambahkan movie_id dan cast_id ke table asosiasi movie_casts
	// melakukan build untuk query insert movie_casts
	queryInsertMovieCasts := utils.InsertTableAssoc("movie_casts", "movie_id", "cast_id", idCasts)

	// mengeksekusi query insert movie_casts yang sudah di build
	if _, err := tx.Exec(ctx, queryInsertMovieCasts, idCasts...); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// tambahkan jadwal untuk movie ini
	// lakukan looping untuk memasukkan jadwal berdasarkan movie yang akan ditampilkan cinema 
	for _, cinema := range cinemaIds {

		// lakukan looping untuk memasukkan jadwal berdasarkan time
		for _, time := range times {
			queryInsertSchedule := "INSERT INTO schedule (cinema_id, movie_id, location, date, time, price) VALUES ($1, $2, $3, $4, $5, $6)"

			if _, err := tx.Exec(ctx, queryInsertSchedule, cinema, movieId, location, date, time, price); err != nil {
				log.Println("[ERROR] : ", err.Error())
				return err
			}

		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}
	return nil
}

// repository update movie (fix)
func (u *MoviesRepository) UpdateMovie(ctx context.Context, title, filePath, overview, directorName string, releaseDate time.Time, duration int, genres, casts []string, idInt int) (pgconn.CommandTag, error) {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	defer tx.Rollback(ctx)

	// melakukan update movie berdasarkan id movie 
	query := "update movies set title = $1, image_path = $2, overview = $3, release_date = $4, director_name = $5, duration = $6, modified_at = $7 where id = $8;"

	values := []any{title, filePath, overview, releaseDate, directorName, duration, time.Now(), idInt}

	cmd, err := tx.Exec(ctx, query, values...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	// memperbaharui movie_genres
	// melakukan delete pada tabel asosiasi movie_genre berdasarkan movie id
	queryDelMovieGenre := "DELETE FROM movie_genres WHERE movie_id = $1"
	_, errDelMovieGenre := tx.Exec(ctx, queryDelMovieGenre, idInt)
	if errDelMovieGenre != nil {
		return pgconn.CommandTag{}, nil
	}
	
	for _, genre := range genres {
		// menambahkan genre baru jika belum terdaftar
		queryGenres := "INSERT INTO genres (name) VALUES ($1) ON CONFLICT (name) DO NOTHING"
		_, err := tx.Exec(ctx, queryGenres, genre)
		if err != nil {
			log.Println("[ERROR] : ", err.Error())
			return pgconn.CommandTag{}, err
		}

		// ambil genre id
		var genreId int
		queryGenreId := "SELECT id FROM genres WHERE name = $1"
        err = tx.QueryRow(ctx, queryGenreId, genre).Scan(&genreId)
        if err != nil {
			log.Println("[ERROR] : ", err.Error())
            return pgconn.CommandTag{}, err
        }

		// tambahkan movie id dan genre id ke tabel asosiasi movie_genre
		queryMovieGenre := "INSERT INTO movie_genres (movie_id, genre_id) VALUES ($1, $2)"
		_, err = tx.Exec(ctx, queryMovieGenre, idInt, genreId)
        if err != nil {
			log.Println("[ERROR] : ", err.Error())
            return pgconn.CommandTag{}, err
        }
	}

		// memperbarui movie_casts
		// melakukan delete pada tabel asosiasi movie_casts berdasarkan movie id
		queryDelMovieCasts := "DELETE FROM movie_casts WHERE movie_id = $1"
		_, errDelMovieCasts := tx.Exec(ctx, queryDelMovieCasts, idInt)
		if errDelMovieCasts != nil {
			log.Println("[ERROR] : ", errDelMovieCasts.Error())
			return pgconn.CommandTag{}, nil
		}

	 	for _, cast := range casts {
		// menambahkan genre baru jika belum terdaftar
		queryGenres := "INSERT INTO casts (name) VALUES ($1) ON CONFLICT (name) DO NOTHING"
		_, err := tx.Exec(ctx, queryGenres, cast)
		if err != nil {
			log.Println("ERROR : ", err.Error())
			return pgconn.CommandTag{}, err
		}

		// ambil genre id
		var castId int
		queryGenreId := "SELECT id FROM casts WHERE name = $1"
        err = tx.QueryRow(ctx, queryGenreId, cast).Scan(&castId)
        if err != nil {
			log.Println("[ERROR]", err.Error())
            return pgconn.CommandTag{}, err
        }

		// tambahkan movie id dan genre id ke tabel asosiasi movie_genre
		queryMovieGenre := "INSERT INTO movie_casts (movie_id, cast_id) VALUES ($1, $2)"
		_, err = tx.Exec(ctx, queryMovieGenre, idInt, castId)
        if err != nil {
			log.Println("[ERROR] : ", err.Error())
            return pgconn.CommandTag{}, err
        }
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}
	

	return cmd, nil
}

// repository delete movie (fix)
func (u *MoviesRepository) DeleteMovie(ctx context.Context, idInt int) (pgconn.CommandTag, error) { 

	tx, err := u.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	defer tx.Rollback(ctx)

	// melakukan delele movie_genres berdasarkan movie_id
	queryMovieGenres := "DELETE FROM movie_genres WHERE movie_id = $1"
	_, errMovieGenres := tx.Exec(ctx, queryMovieGenres, idInt)
	if errMovieGenres != nil {
		log.Println("[ERROR] : ", errMovieGenres.Error())
		return pgconn.CommandTag{}, nil
	}

	// melakukan delete movie_casts berdasarkan movie_id
	queryMovieCasts := "DELETE FROM movie_casts WHERE movie_id = $1"
	_, errMovieCasts := tx.Exec(ctx, queryMovieCasts, idInt)
	if errMovieCasts != nil {
		log.Println("[ERROR] : ", errMovieCasts.Error())
		return pgconn.CommandTag{}, nil
	}

	// melakukan delete schedule berdasarkan movie_id
	querySchedule := "DELETE FROM schedule WHERE movie_id = $1"
	_, errSchedule := tx.Exec(ctx, querySchedule, idInt)
	if errSchedule != nil {
		log.Println("[ERROR] : ", errSchedule.Error())
		return pgconn.CommandTag{}, nil
	}

	// melakukan delete movie berdasarkan movie id
	queryMovie := "DELETE FROM movies WHERE id = $1"
	cmd, err := tx.Exec(ctx, queryMovie, idInt)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, nil
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	return cmd, nil
}

// repository get upcoming movie (fix)
func (u *MoviesRepository) GetMovieUpcoming(ctx context.Context) ([]models.MoviesStruct, error) {
	
	// mengambil data movies dengan ketentuan release_date harus lebih besar dari pada waktu sekarang
	query := `select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(distinct g.name) as "genres", array_agg(distinct c.name) from movies m join movie_genres mg on mg.movie_id = m.id join genres g on mg.genre_id = g.id join movie_casts mc on mc.movie_id = m.id join casts c on c.id = mc.cast_id where m.release_date > now() group by m.id order by m.id asc`
	rows, err := u.db.Query(ctx, query)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
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

// repository get popular movie
func (u *MoviesRepository) GetMoviePopular(ctx context.Context) ([]models.MoviesStruct, error) {

	// mengambil data movies dengan ketentuan jumlah jumlah order dalam tabel orders_seat harus lebih dari 10
	query := `select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(distinct g.name) as "genres", array_agg(distinct c.name) as "casts", COUNT(distinct os.id) as "qty" from orders o join schedule s on o.schedule_id = s.id join movies m on m.id = s.movie_id join order_seats os on os.order_id = o.id join movie_genres mg on mg.movie_id = m.id join genres g on g.id = mg.genre_id join movie_casts mc on mc.movie_id = m.id join casts c on c.id = mc.cast_id group by m.id having COUNT(os.order_id) > 10`

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_movie, &movies.Duration, &movies.Director_name, &movies.Genres, &movies.Casts, &movies.TotalSales); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	return result, nil
}

// repository get detail movie (fix)
func (u *MoviesRepository) GetDetailMovie(ctx context.Context, movies models.MoviesStruct, IdInt int) ([]models.MoviesStruct, error) {

	// menambil data movies berdasarkan movie_id
	query := `with table_movie_genres as (select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(g.name) as "genres" from movies m join movie_genres mg on m.id = mg.movie_id join genres g on g.id = mg.genre_id where m.id = $1 group by m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name) select t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres, array_agg(c.name) from table_movie_genres t join movie_casts mc on t.id = mc.movie_id join casts c on c.id = mc.cast_id group by t.id, t.title, t.release_date, t.overview, t.image_path, t.duration, t.director_name, t.genres`
	// values := []any{IdInt}
	rows, err := u.db.Query(ctx, query, IdInt)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
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

// repository get movie with pagination (fix)
func (u *MoviesRepository) GetMoviesWithPagination(ctx context.Context, movie models.MoviesStruct, offset int, title, genre string) ([]models.MoviesStruct, error) {

	// mengambil data movies didalam subquery yang sudah dibatasi menggunakan limit dan offset, dan hasil data movie tersebut dilakukan filter berdasarkan title dan genre movienya
	query := `select id, title, release_date, overview, image_path, duration, director_name, genres, casts from (select m.id, m.title, m.release_date, m.overview, m.image_path, m.duration, m.director_name, array_agg(distinct g.name) as "genres", array_agg(distinct c.name) as "casts" from movies m join movie_genres mg on mg.movie_id = m.id join genres g on mg.genre_id = g.id join movie_casts mc on mc.movie_id = m.id join casts c on c.id = mc.cast_id group by m.id order by m.id limit 5 offset $1) sq where lower(sq.title) like '%' || lower($2) ||'%' and lower(array_to_string(sq.genres, ',')) like '%' || lower($3) || '%'`
	values := []any{offset, title, genre}
	rows, err := u.db.Query(ctx, query, values...)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}

	defer rows.Close()
	var result []models.MoviesStruct
	for rows.Next() {
		var movies models.MoviesStruct
		if err := rows.Scan(&movies.Id, &movies.Title, &movies.Release_date, &movies.Overview, &movies.Image_movie, &movies.Duration, &movies.Director_name, &movies.Genres, &movies.Casts); err != nil {
			return nil, err
		}
		result = append(result, movies)
	}
	
	return result, nil
}