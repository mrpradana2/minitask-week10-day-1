package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
    
    // DATA USERS
    type usersStruct struct {
        Id int `json:"id" form:"id"`
        Email string `json:"email" form:"email"`
        Password string `json:"password" form:"password"`
        Role string `json:"role" form:"role"`
    }

    users := []usersStruct{
        {Id: 1, Email: "pradana@gmail.com", Password: "pradana123", Role: "user"},
        {Id: 2, Email: "jhon@gmail.com", Password: "jhon123", Role: "user"},
        {Id: 3, Email: "dhani@gmail.com", Password: "dhani123", Role: "user"},
    }

    router.GET("/users", func(ctx *gin.Context) {
        emailQ := ctx.Query("Email")
        passQ := ctx.Query("Password")

        // jika tidak terdapat query apapun
        if emailQ == "" && passQ == "" {
            ctx.JSON(http.StatusOK, gin.H{
                "msg": "success",
                "data": users,
            })
            return
        }

        // mencari user dan password
        result := []usersStruct{}
        errPass := ""
        for _, user := range users {
            emailCond := strings.EqualFold(user.Email, emailQ)
            passCond := strings.EqualFold(user.Password, passQ)

            if emailCond && passCond {
                result = append(result, user)
            }

            if emailCond && !passCond {
                errPass = "incorrect password"
            }
        }

        // jika password salah
        if errPass != "" {
            ctx.JSON(http.StatusNotFound, gin.H{
                "msg": "incorrect password",
            })
            return
        }

        // jika user tidak ditemukan
        if len(result) == 0 {
            ctx.JSON(http.StatusNotFound, gin.H{
                "msg": "user not found",
            })
            return
        }

        // jika user ditemukan
        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": result,
        })
    })

    // menambahkan user baru
    router.POST("users", func(ctx *gin.Context) {
        newUser := &usersStruct{}

        if err := ctx.ShouldBind(newUser); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "msg": "an error occurred on the server",
            })
            return
        }

        newUsers := append(users, *newUser)
        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": newUsers,
        })
    })

    // DATA PROFILE
    type profileStruct struct {
        UserId int `json:"userId" form:"userId"`
        FirstName string `json:"firstName" form:"firstName"`
        LastName string `json:"lastName" form:"lastName"`
        PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`
        PhotoPath string `json:"photoPath" form:"photoPath"`
        Title string `json:"title" form:"title"`
        Point int `json:"point" form:"point"`
    }

    profiles := []profileStruct{
        {UserId: 1, FirstName: "Rizki", LastName: "Pradana", PhoneNumber: "087768012234", PhotoPath: "/images/photo/pradana.png", Title: "Director of Finance", Point: 0},
        {UserId: 2, FirstName: "Doni", LastName: "Rahmawan", PhoneNumber: "087768010011", PhotoPath: "/images/photo/rahmawan.png", Title: "Teacher", Point: 0},
        {UserId: 3, FirstName: "Tony", LastName: "Sugiharto", PhoneNumber: "087768018822", PhotoPath: "/images/photo/sugiharto.png", Title: "Director", Point: 0},
        {UserId: 4, FirstName: "Tony", LastName: "Sugiharto", PhoneNumber: "087768018822", PhotoPath: "/images/photo/sugiharto.png", Title: "Director", Point: 0},
        {UserId: 5, FirstName: "Fany", LastName: "Rahmawati", PhoneNumber: "087768012221", PhotoPath: "/images/photo/rahmawati.png", Title: "Teller Bank", Point: 0},
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

        var profile []profileStruct
        for _, p := range profiles {
            if p.UserId == idInt {
                profile = append(profile, p)
                break
            }
        } 

        if len(profile) == 0 {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "msg": "profile not found",
            })
            return
        }

        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": profile[0],
        })

    })

    router.GET("/profiles", func(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, gin.H{
            "msg": "success",
            "data": profiles,
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
            "msg": "success",
            "data": movies,
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