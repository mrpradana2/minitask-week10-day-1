package handlers

import (
	"fmt"
	"log"
	"net/http"
	fp "path/filepath"
	"strconv"
	"strings"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"
	"tikcitz-app/pkg"
	"time"

	"github.com/gin-gonic/gin"
)

type Movieshandler struct{
	moviesRepo *repositories.MoviesRepository
}

func NewMovieshandler(moviesRepo *repositories.MoviesRepository) *Movieshandler {
	return &Movieshandler{moviesRepo: moviesRepo}
}

// Handler get all movies (fix)
func (m *Movieshandler) GetMovies(ctx *gin.Context) {

	// manjalankan fungsi repository get movies 
	result, err := m.moviesRepo.GetMovies(ctx.Request.Context())

	// mengecek jika terjadi error saat mengakses data di server
	if err != nil {
		log.Println("[ERROR]", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// mengecek jika data yang diambil dari server kosong
	if len(result) == 0 {
		log.Println(result)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "movie not found",
		})
		return
	}

	// menampilkan hasil response jika request berhasil dari server
	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "success",
		Result: result,
	})
}

// Handler add movie (fix)
func (m *Movieshandler) AddMovie(ctx *gin.Context)  {
	// siapkan variable movie struct
	var newDataMovie models.MoviesStruct

	// lakukan binding data jika ada data yang tidak sesuai dengan format, maka akan terjadi error
	if err := ctx.ShouldBind(&newDataMovie); err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "Field validation") {
			ctx.JSON(http.StatusInternalServerError, models.Message{
				Status: "failed",
				Msg: "Ada content yang harus diisi",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "terjadi kesalahan serverr",
		})
		return
	}

	// ambil masing masing new data dari form
	file := newDataMovie.Image_path
	title := newDataMovie.Title
	overview := newDataMovie.Overview
	releaseDate := newDataMovie.Release_date
	directorName := newDataMovie.Director_name
	duration := newDataMovie.Duration
	genres := newDataMovie.Genres
	casts := newDataMovie.Casts
	cinemaIds := newDataMovie.Cinema_ids
	location := newDataMovie.Location
	date := newDataMovie.Date
	times := newDataMovie.Times
	price := newDataMovie.Price

	log.Println("CASTS : ", casts)
	log.Println("GENRES : ", genres)
	
	// jika file bernilai nil maka tampilkan error
	if file == nil {
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "your file is empty",
		})
		return
	}

	// ambil id yang ada di header
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)
	ext := fp.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d_movie_image%s", time.Now().UnixNano(), userClaims.Id, ext)
	filepath := fp.Join("public", "img", filename)
	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "terjadi kesalahan upload",
		})
		return
	}

	// jalankan fungsi untuk menambahkan data movie
	err := m.moviesRepo.AddMovie(ctx.Request.Context(), title, filepath, overview, directorName, location, releaseDate, date, times, duration, price, genres, casts, cinemaIds)

	// error jika terjadi masalah saat mengirim data
	if err != nil {
		log.Println("[ERROR]:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// tampilkan pesan success jika data berhasil ditambahkan
	ctx.JSON(http.StatusCreated, models.Message{
		Status: "ok",
		Msg: "successfully add data movie",
	})
}

// Handler update movie (fix)
func (m *Movieshandler) UpdateMovie(ctx *gin.Context) {

	// mengambil id dari params
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "error",
			Msg: "Param id is needed",
		})
		return
	}

	// melakukan konversi id string menjadi id integer
	idInt, err := strconv.Atoi(idStr)

	// handling error jika gagal mengkonversi id
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// siapkan variable movie struct
	var newDataMovie models.MoviesStruct

	// lakukan binding data jika ada data yang tidak sesuai dengan format, maka akan terjadi error
	if err := ctx.ShouldBind(&newDataMovie); err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "Field validation") {
			ctx.JSON(http.StatusInternalServerError, models.Message{
				Status: "failed",
				Msg: "Ada content yang harus diisi",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "terjadi kesalahan serverr",
		})
		return
	}

	// ambil masing masing new data dari form
	file := newDataMovie.Image_path
	title := newDataMovie.Title
	overview := newDataMovie.Overview
	releaseDate := newDataMovie.Release_date
	directorName := newDataMovie.Director_name
	duration := newDataMovie.Duration
	genres := newDataMovie.Genres
	casts := newDataMovie.Casts
	log.Println("hello", genres)
	// jika file bernilai nil maka tampilkan error
	if file == nil {
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "your file is empty",
		})
		return
	}

	// ambil id yang ada di header
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)
	ext := fp.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d_movie_image%s", time.Now().UnixNano(), userClaims.Id, ext)
	filepath := fp.Join("public", "img", filename)
	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "failed",
			Msg: "terjadi kesalahan upload",
		})
		return
	}

	// menjalankan fungsi repository update movie
	cmd, err := m.moviesRepo.UpdateMovie(ctx.Request.Context(), title, filepath, overview, directorName, releaseDate, duration, genres, casts, idInt)

	// melakukan handling error jika query gagal
	if err != nil {
		log.Println("Insert profile error:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// handling error jika tidak terjadi perubahan sama sekali di database
	if cmd.RowsAffected() == 0 {
		log.Println("Query failed, could not change the data in the database")
	}

	// menampilkan response dari server, jika request dari client berhasil
	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "data successfully changed",
	})
}

// handler delete movie (fix)
func (m *Movieshandler) DeleteMovie(ctx *gin.Context) {
	// mendapatkan id dari params
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "error",
			Msg: "Param id is needed",
		})
		return
	}

	// melakukan konversi id string menjadi id integer
	idInt, err := strconv.Atoi(idStr)

	// melakukan error handling jika gagal mengkonversi id
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// menjalankan fungsi repository delete movie 
	cmd, err := m.moviesRepo.DeleteMovie(ctx.Request.Context(), idInt)

	// melakukan error handling jika query gagal dijalankan
	if err != nil {
		log.Println("[ERROR]", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// error handling jika movie yang akan dihapus tidak tersedia
	if cmd.RowsAffected() == 0 {
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: "failed",
			Msg: "Movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "data successfully deleted",
	})
}

// handler get movie upcoming  (fix)
func (m *Movieshandler) GetMoviesUpcoming(ctx *gin.Context) {

	// menjalankan fungsi repository get movie upcoming
	result, err := m.moviesRepo.GetMovieUpcoming(ctx.Request.Context())

	// mengecek jika terjadi error saat mengakses data di server
	if err != nil {
		log.Println("Get Movie error:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "ok",
			Msg: "an error occurred on the server",
		})
		return
	}

	// error handling jika hasil request dari server kosong
	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: "ok",
			Msg: "movie not found",
		})
	}

	// tampilkan hasil jika berhasil menerima request dari server
	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "success",
		Result: result,
	})
}

// handler get movie popular
func (m *Movieshandler) GetMoviesPopular(ctx *gin.Context) {
	result, err := m.moviesRepo.GetMoviePopular(ctx.Request.Context())
	if err != nil {
		log.Println("Get Movie error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "movie not found",
			"data": result,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": result,
	})
}

// handler get detail movie (fix)
func (m *Movieshandler) GetDetailMovie(ctx *gin.Context) {
	// ambil id param
	idStr, ok := ctx.Params.Get("id")

	// handling error jika param tidak ada
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "error",
			Msg: "Param id is needed",
		})
		return
	}

	// konversi id string menjadi id integer
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		log.Println("[BEDUG] : ", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// jalankan fungsi repository untuk get detail movie
	result, err := m.moviesRepo.GetDetailMovie(ctx.Request.Context(), models.MoviesStruct{}, idInt)
	if err != nil {
		log.Println("Get Movie error:", err)
		ctx.JSON(http.StatusInternalServerError, models.Message{
			Status: "error",
			Msg: "an error occurred on the server",
		})
		return
	}

	// mengecek jika request kosong
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, models.Message{
			Status: "not found",
			Msg: "movie not found",
		})
		return
	}

	// menampilkan hasil response server jika request client berhasil 
	ctx.JSON(http.StatusOK, models.Message{
		Status: "ok",
		Msg: "success",
		Result: result,
	})
}

// handler get movie with pagination
func (m *Movieshandler) GetMoviesWithPagination(ctx *gin.Context) {
	pageQ := ctx.Query("page")

	if pageQ == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "query page is needed",
		})
		return
	}

	pageQInt, err := strconv.Atoi(pageQ)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "an error occurred on the server",
		})
		return
	}

	var offset int
	if pageQInt == 1 {
		offset = 0
	} else {
		offset = pageQInt * 5 - 5
	}

	result, err := m.moviesRepo.GetMoviesWithPagination(ctx.Request.Context(), models.MoviesStruct{}, offset)

	if err != nil {
		log.Println("Get Movie error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": result,
	})
}