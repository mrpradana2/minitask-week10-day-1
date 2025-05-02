package handlers

import (
	"log"
	"net/http"
	"strconv"
	"tikcitz-app/internals/models"
	"tikcitz-app/internals/repositories"

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
	// mengambil body json / input admin
	newDataMovie := models.MoviesStruct{}

	// binding data 
	// mambaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&newDataMovie); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "error",
			Msg: "invalid data sent",
		})
		return
	}

	// jalankan fungsi untuk menambahkan data movie
	err := m.moviesRepo.AddMovie(ctx.Request.Context(), newDataMovie)

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

// Handler update movie
func (m *Movieshandler) UpdateMovie(ctx *gin.Context) {
	var updateMovie models.MoviesStruct

	// binding data
	// mambaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error 
	if err := ctx.ShouldBindJSON(&updateMovie); err != nil {
		log.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, models.Message{
			Status: "error",
			Msg: "invalid data sent",
		})
		return
	}

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

	// menjalankan fungsi repository update movie
	cmd, err := m.moviesRepo.UpdateMovie(ctx.Request.Context(), &updateMovie, idInt)

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

// handler delete movie
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

// handler get movie upcoming
func (m *Movieshandler) GetMoviesUpcoming(ctx *gin.Context) {
	result, err := m.moviesRepo.GetMovieUpcoming(ctx.Request.Context())
	if err != nil {
		log.Println("Get Movie error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "an error occurred on the server"})
		return
	}

	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "movie not found",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": result,
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