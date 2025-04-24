package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := gin.Default()

	// setting file env yang berisi informasi db
    dbEnv := []any{}
	dbEnv = append(dbEnv, os.Getenv("DBUSER"))
	dbEnv = append(dbEnv, os.Getenv("DBPASS"))
	dbEnv = append(dbEnv, os.Getenv("DBHOST"))
	dbEnv = append(dbEnv, os.Getenv("DBPORT"))
	dbEnv = append(dbEnv, os.Getenv("DBNAME"))

	// setup database connection
	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbEnv...)

	dbClient, err := pgxpool.New(context.Background(), dbString)

	log.Println("[debug]", dbString)

	if err != nil {
		log.Printf("unable to create connection pool:  %v\n", err)
		os.Exit(1)
	}

	// closing DB
	defer func()  {
		log.Println("Closing DB...")
		dbClient.Close()	
	}()

    // DATA USER
    type usersStruct struct {
        Id int `db:"id" json:"id.omitempty"`
        Email string `db:"email" json:"email" form:"email"`
        Password string `db:"password" json:"password" form:"password"`
    }

    // add user
    router.POST("/add-users", func(ctx *gin.Context) {
        newUser := usersStruct{}

        if err := ctx.ShouldBindJSON(&newUser); err != nil {
			log.Println("Binding error:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "invalid data sent",
			})
			return
		}

        query := "INSERT INTO users (email, password) VALUES ($1, $2)"
		values := []any{newUser.Email, newUser.Password}
		cmd, err := dbClient.Exec(ctx.Request.Context(), query, values...)
        if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Terjadi kesalahan server",
			})
			return
		}

		if cmd.RowsAffected() == 0 {
			log.Println("Query gagal, tidak merubah data di DB")
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
    })

    // auth user login
    router.POST("/user-login", func(ctx *gin.Context) {
        auth := usersStruct{}

		if err := ctx.ShouldBindJSON(&auth); err != nil {
			log.Println("Binding error:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "Data yang dikirim tidak valid",
			})
			return
		}

        // mengambil data user dari DB
        query := "SELECT email, password FROM users"
        values := []any{auth.Email, auth.Password}

        rows, err := dbClient.Query(context.Background(), query)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "terjadi kesalahan sistem",
			})
			return
		}

		defer rows.Close()
		var result []usersStruct
		for rows.Next() {
			var users usersStruct
			if err := rows.Scan(&users.Email, &users.Password); err != nil {
				log.Println(err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "terjadi kesalahan sistem",
				})
				return
			}
			result = append(result, users)
		}

        // mengecek apakan email dan password sesuai
        var userLogin []usersStruct
        for _, user := range result {
            if user.Email == values[0] && user.Password == values[1] {
                userLogin = append(userLogin, user)
            }
        }

        // error handling jika email atau password tidak sesuai
        if len(userLogin) == 0 {
            ctx.JSON(http.StatusNotFound, gin.H{
                "msg": "incorrect email or password",
            })
            return
        }

        // jika user berhasil login
        ctx.JSON(http.StatusOK, gin.H{
            "msg": "login user success",
        })
    })

    // DATA PROFILE
    type profileStruct struct {
        UserId int `db:"userId" json:"userId" form:"userId.omitempty"`
        FirstName string `db:"firstName" json:"firstName" form:"firstName"`
        LastName string `db:"lastName" json:"lastName" form:"lastName"`
        PhoneNumber string `db:"phoneNumber" json:"phoneNumber" form:"phoneNumber"`
        PhotoPath string `db:"photoPath" json:"photoPath" form:"photoPath"`
        Title string `db:"title" json:"title" form:"title"`
        Point int `db:"point" json:"point" form:"point.omitempty"`
    }


    router.GET("/profiles/:id", func(ctx *gin.Context) {
        idStr, ok := ctx.Params.Get("id")

        // handling error jika param tidak ada
        if !ok {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "msg": "Param id is needed",
            })
            return
        }

        idInt, err := strconv.Atoi(idStr)

        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "msg": "an error occurred on the server",
            })
            return
        }

        query := "SELECT phone_number, first_name, last_name, photo_path, title, point FROM profile WHERE user_id = $1"
        values := []any{idInt}
        var result profileStruct
		if err := dbClient.QueryRow(context.Background(), query, values...).Scan(&result.PhoneNumber, &result.FirstName, &result.LastName, &result.PhotoPath, &result.Title, &result.Point); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "terjadi kesalahan sistem",
			})
			return
		}

        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": result,
        })
    })

    // DATA MOVIES
    type moviesStruct struct {
        Id int `json:"id" form:"id"`
        Title string `json:"title" form:"title"`
        Image_path string `json:"image_path" form:"image_path"`
        Overview string `json:"overview" form:"overview"`
        Release_date string `json:"release_date" form:"release_date"`
        Director_name string `json:"director_name" form:"director_name"`
        Duration int `json:"duration" form:"duration"`
        Casts []string `json:"casts" form:"casts"`
        Status_movie_id int `json:"status_movie_id" form:"status_movie_id"`
        Genres []string `json:"genres" form:"genres"`
        Status_movie string `json:"status_movie" form:"status_movie"`
    }

    var movies = []moviesStruct{
        {
            Id: 1,
            Title: "Avengers: Secret Wars",
            Image_path: "https://example.com/images/avengers_secret_wars.jpg",
            Overview: "Earth's mightiest heroes reunite to face a multiversal threat.",
            Release_date: "2026-05-01",
            Director_name: "Destin Daniel Cretton",
            Duration: 150,
            Casts: []string{"Robert Downey Jr.", "Chris Evans", "Scarlett Johansson"},
            Status_movie_id: 1,
            Genres: []string{"Action", "Adventure", "Sci-Fi"},
            Status_movie: "movie upcoming",
        },
        {
            Id: 2,
            Title: "Dune: Part Two",
            Image_path: "https://example.com/images/dune_part_two.jpg",
            Overview: "Paul Atreides unites with the Fremen to fulfill his destiny.",
            Release_date: "2024-11-20",
            Director_name: "Denis Villeneuve",
            Duration: 165,
            Casts: []string{"Timothée Chalamet", "Zendaya", "Rebecca Ferguson"},
            Status_movie_id: 2,
            Genres: []string{"Adventure", "Drama", "Sci-Fi"},
            Status_movie: "movie popular",
        },
        {
            Id: 3,
            Title: "The Batman: Part II",
            Image_path: "https://example.com/images/batman_2.jpg",
            Overview: "Batman faces a new wave of crime in Gotham City.",
            Release_date: "2025-10-03",
            Director_name: "Matt Reeves",
            Duration: 155,
            Casts: []string{"Robert Pattinson", "Zoë Kravitz", "Jeffrey Wright"},
            Status_movie_id: 1,
            Genres: []string{"Action", "Crime", "Drama"},
            Status_movie: "movie upcoming",
        },
        {
            Id: 4,
            Title: "Spider-Man: Beyond the Spider-Verse",
            Image_path: "https://example.com/images/spiderverse.jpg",
            Overview: "Miles Morales continues his journey across the multiverse.",
            Release_date: "2025-03-15",
            Director_name: "Joaquim Dos Santos",
            Duration: 130,
            Casts: []string{"Shameik Moore", "Hailee Steinfeld", "Oscar Isaac"},
            Status_movie_id: 1,
            Genres: []string{"Animation", "Action", "Adventure"},
            Status_movie: "movie upcoming",
        },
        {
            Id: 5,
            Title: "Mission: Impossible – Dead Reckoning Part Two",
            Image_path: "https://example.com/images/mission_impossible_8.jpg",
            Overview: "Ethan Hunt returns for another impossible mission.",
            Release_date: "2025-06-27",
            Director_name: "Christopher McQuarrie",
            Duration: 160,
            Casts: []string{"Tom Cruise", "Ving Rhames", "Simon Pegg"},
            Status_movie_id: 1,
            Genres: []string{"Action", "Thriller"},
            Status_movie: "movie upcoming",
        },
        {
            Id: 6,
            Title: "Oppenheimer",
            Image_path: "https://example.com/images/oppenheimer.jpg",
            Overview: "The story of J. Robert Oppenheimer and the atomic bomb.",
            Release_date: "2023-07-21",
            Director_name: "Christopher Nolan",
            Duration: 180,
            Casts: []string{"Cillian Murphy", "Emily Blunt", "Matt Damon"},
            Status_movie_id: 2,
            Genres: []string{"Biography", "Drama", "History"},
            Status_movie: "movie popular",
        },
        {
            Id: 7,
            Title: "The Marvels",
            Image_path: "https://example.com/images/the_marvels.jpg",
            Overview: "Captain Marvel joins forces with Ms. Marvel and Monica Rambeau.",
            Release_date: "2023-11-10",
            Director_name: "Nia DaCosta",
            Duration: 105,
            Casts: []string{"Brie Larson", "Iman Vellani", "Teyonah Parris"},
            Status_movie_id: 2,
            Genres: []string{"Action", "Adventure", "Fantasy"},
            Status_movie: "movie popular",
        },
        {
            Id: 8,
            Title: "Wonka",
            Image_path: "https://example.com/images/wonka.jpg",
            Overview: "A young Willy Wonka embarks on a magical adventure.",
            Release_date: "2023-12-15",
            Director_name: "Paul King",
            Duration: 116,
            Casts: []string{"Timothée Chalamet", "Olivia Colman", "Hugh Grant"},
            Status_movie_id: 2,
            Genres: []string{"Adventure", "Comedy", "Family"},
            Status_movie: "movie popular",
        },
        {
            Id: 9,
            Title: "The Hunger Games: The Ballad of Songbirds & Snakes",
            Image_path: "https://example.com/images/hunger_games_prequel.jpg",
            Overview: "A prequel to the Hunger Games saga, set decades before Katniss.",
            Release_date: "2023-11-17",
            Director_name: "Francis Lawrence",
            Duration: 157,
            Casts: []string{"Tom Blyth", "Rachel Zegler", "Hunter Schafer"},
            Status_movie_id: 2,
            Genres: []string{"Action", "Drama", "Sci-Fi"},
            Status_movie: "movie popular",
        },
        {
            Id: 10,
            Title: "Aquaman and the Lost Kingdom",
            Image_path: "https://example.com/images/aquaman2.jpg",
            Overview: "Aquaman forges an uneasy alliance to save Atlantis.",
            Release_date: "2023-12-20",
            Director_name: "James Wan",
            Duration: 124,
            Casts: []string{"Jason Momoa", "Patrick Wilson", "Amber Heard"},
            Status_movie_id: 2,
            Genres: []string{"Action", "Adventure", "Fantasy"},
            Status_movie: "movie popular",
        },
        {
            Id: 11,
            Title: "Inside Out 2",
            Image_path: "https://example.com/images/inside_out_2.jpg",
            Overview: "Riley enters her teenage years with brand-new emotions.",
            Release_date: "2025-06-13",
            Director_name: "Kelsey Mann",
            Duration: 95,
            Casts: []string{"Amy Poehler", "Phyllis Smith", "Maya Hawke"},
            Status_movie_id: 1,
            Genres: []string{"Animation", "Comedy", "Family"},
            Status_movie: "movie upcoming",
        },
        {
            Id: 12,
            Title: "Fantastic Four",
            Image_path: "https://example.com/images/fantastic_four.jpg",
            Overview: "Marvel's First Family enters the MCU.",
            Release_date: "2025-07-25",
            Director_name: "Matt Shakman",
            Duration: 135,
            Casts: []string{"Pedro Pascal", "Vanessa Kirby", "Joseph Quinn"},
            Status_movie_id: 1,
            Genres: []string{"Action", "Adventure", "Sci-Fi"},
            Status_movie: "movie upcoming",
        },
        {
            Id: 13,
            Title: "Deadpool & Wolverine",
            Image_path: "https://example.com/images/deadpool_wolverine.jpg",
            Overview: "Deadpool joins forces with Wolverine in a multiverse adventure.",
            Release_date: "2024-07-26",
            Director_name: "Shawn Levy",
            Duration: 115,
            Casts: []string{"Ryan Reynolds", "Hugh Jackman", "Morena Baccarin"},
            Status_movie_id: 1,
            Genres: []string{"Action", "Comedy", "Sci-Fi"},
            Status_movie: "movie upcoming",
        },
        {
            Id: 14,
            Title: "Guardians of the Galaxy Vol. 3",
            Image_path: "https://example.com/images/guardians3.jpg",
            Overview: "The Guardians face the aftermath of Gamora's loss.",
            Release_date: "2023-05-05",
            Director_name: "James Gunn",
            Duration: 149,
            Casts: []string{"Chris Pratt", "Zoe Saldana", "Dave Bautista"},
            Status_movie_id: 2,
            Genres: []string{"Action", "Adventure", "Comedy"},
            Status_movie: "movie popular",
        },
        {
            Id: 15,
            Title: "The Flash",
            Image_path: "https://example.com/images/flash.jpg",
            Overview: "Barry Allen travels back in time to save his mother.",
            Release_date: "2023-06-16",
            Director_name: "Andy Muschietti",
            Duration: 144,
            Casts: []string{"Ezra Miller", "Michael Keaton", "Sasha Calle"},
            Status_movie_id: 2,
            Genres: []string{"Action", "Adventure", "Fantasy"},
            Status_movie: "movie popular",
        },
    }

    router.GET("/movies", func(ctx *gin.Context) {

		if len(movies) == 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "movie not found",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"data": movies,
		})
	})

    router.POST("movies", func(ctx *gin.Context) {
        newMovie := &moviesStruct{}
        if err := ctx.ShouldBind(newMovie); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "msg": "terjadi kesalahan sisten",
            })
            return
        }

        newMovies := append(movies, *newMovie)
        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": newMovies,
        })
    })

    router.GET("movies/movie-upcoming", func(ctx *gin.Context) {
        result := []moviesStruct{}

        if len(movies) == 0 {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "msg": "an error occurred on the server",
            })
        }

        for _, movie := range movies {
            if movie.Status_movie == "movie upcoming" {
                result = append(result, movie)
            }
        }

        if len(result) == 0 {
            ctx.JSON(http.StatusNotFound, gin.H{
                "msg": "movie upcoming not found",
            })
        }

        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": result,
        })
    })
    
    router.GET("movies/movie-popular", func(ctx *gin.Context) {
        result := []moviesStruct{}

        if len(movies) == 0 {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "msg": "an error occurred on the server",
            })
        }

        for _, movie := range movies {
            if movie.Status_movie == "movie popular" {
                result = append(result, movie)
            }
        }

        if len(result) == 0 {
            ctx.JSON(http.StatusNotFound, gin.H{
                "msg": "movie popular not found",
            })
        }

        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": result,
        })
    })

    router.Run()
}