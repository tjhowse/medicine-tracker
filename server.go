package main

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./web/api-codegen-config.yaml ./web/api.yaml

type Server struct {
	db       *MedicineLoggerDB
	e        *echo.Echo
	settings ServerSettings
	email    *email
}

type ServerSettings struct {
	DBString        string `default:"medicine-logger.sqlite3"`
	JWTSecret       string `default:"very_very_secret"`
	Address         string `default:":3019"`
	AutoTLS         bool   `default:"false"`
	AccountCreation bool   `default:"true"`
	EmailHost       string `default:"smtp.sendgrid.net"`
	EmailPort       int    `default:"587"`
	EmailUsername   string `default:"apikey"`
	EmailPassword   string `default:""`
	EmailFrom       string `default:"friend@medicine-logger.com"`
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (s *Server) Init(settings ServerSettings) {
	var err error
	s.settings = settings
	if s.db, err = NewMedicineLoggerDB(s.settings.DBString); err != nil {
		log.Fatal(err)
	}

	s.email = GetEmail(&s.settings)

	s.e = echo.New()

	// s.e.Use(middleware.Logger())
	// s.e.Use(middleware.Recover())
	// s.e.Logger.SetLevel(2)

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:  []byte(s.settings.JWTSecret),
		TokenLookup: "cookie:token",
	}

	s.e.GET("/login.html", s.login)

	// Redirect requests missing a JWT token to the login page.
	jwtConfig.ErrorHandler = func(c echo.Context, err error) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/login.html")
	}

	s.e.Static("/", "web")
	// TODO fix this prefix so the effective endpoint isn't /api/v1/api/v1/...
	api := s.e.Group("/api/v1")
	api.Use(echojwt.WithConfig(jwtConfig))

	mainPage := s.e.Group("/index.html")
	mainPage.Use(echojwt.WithConfig(jwtConfig))
	RegisterHandlersWithBaseURL(api, s, "")
}

func (s *Server) Run() {
	s.e.Logger.Fatal(s.e.Start(s.settings.Address))
}

func (s *Server) Close() {
	s.db.Close()
}

func (s *Server) login(ctx echo.Context) error {
	// TODO write an API entry for this endpoint.
	username := ctx.FormValue("username")
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	oneTimeToken := ctx.FormValue("ott")

	var claims *jwtCustomClaims

	// If we have a one-time token, use it to log in.
	if oneTimeToken != "" {
		if username, err := s.db.ValidateOneTimeToken(oneTimeToken); err != nil {
			log.Println("Invalid one-time token: ", err)
			return echo.ErrInternalServerError
		} else if username == "" {
			log.Println("Invalid one-time token")
			return echo.ErrUnauthorized
		} else {
			// Set custom claims
			claims = &jwtCustomClaims{
				string(username),
				false,
				jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
				},
			}
		}
	}

	// If no form data, serve login.html
	// TODO this has the strange effect of serving login.html to unauthenticated
	// calls to API endpoints.
	if username == "" || password == "" {
		return ctx.File("web/login.html")
	}

	// Throws unauthorized error
	if extant, valid, err := s.db.ValidateUser(username, password); err != nil {
		log.Println("Invalid username/password: ", err)
		return echo.ErrInternalServerError
	} else if extant && !valid {
		log.Println("Invalid password for user ", username)
		return echo.ErrUnauthorized
	} else if !extant {
		if s.settings.AccountCreation {
			// Create a new user
			if err := s.db.AddUser(username, email, password); err != nil {
				log.Println("Failed to add user: ", err)
				return echo.ErrInternalServerError
			} else {
				log.Println("Created new user", username)
			}
		} else {
			log.Println("User ", username, " does not exist")
			return echo.ErrUnauthorized
		}
	}

	// Set custom claims
	claims = &jwtCustomClaims{
		username,
		false,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(s.settings.JWTSecret))
	if err != nil {
		return err
	}
	// Set the JWT as a cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour * 24 * 365)
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	ctx.SetCookie(cookie)
	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

// Extract the user GUID from the JWT cookie, or return an empty string if the
// cookie is missing or invalid.
func (s *Server) GetOrCreateUserGUIDCookie(ctx echo.Context) (userGUID, error) {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		return userGUID(""), echo.ErrUnauthorized
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.settings.JWTSecret), nil
	})
	if err != nil {
		log.Println("Failed to parse token: ", err)
		return userGUID(""), echo.ErrUnauthorized
	}
	claims := token.Claims.(*jwtCustomClaims)

	// Check that username is in the database
	if exists, _, err := s.db.ValidateUser(claims.Name, ""); err != nil || !exists {
		return userGUID(""), echo.ErrUnauthorized
	}
	return userGUID(claims.Name), nil

}

func (s *Server) GetApiV1MedicineLog(ctx echo.Context, params GetApiV1MedicineLogParams) error {
	var user userGUID
	var err error
	if user, err = s.GetOrCreateUserGUIDCookie(ctx); err != nil {
		return err
	}
	var start, end time.Time
	log.Println("Got a request for the medicine log")
	if params.Start == nil {
		// No start time specified, use a point in the distant past
		start = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		start = *params.Start
	}
	if params.End == nil {
		// No end time specified, use the current time
		end = time.Now()
	} else {
		end = *params.End
	}
	if params.Start != nil || params.End != nil {
		// We don't support filtering based on time, so return an error.
		return echo.ErrBadRequest
	} else {
		if result, err := s.db.GetMedicineLog(user, start, end); err != nil {
			return err
		} else {
			return ctx.JSON(http.StatusOK, result)
		}
	}
}

func (s *Server) PostApiV1MedicineLog(ctx echo.Context) error {
	var user userGUID
	var err error
	if user, err = s.GetOrCreateUserGUIDCookie(ctx); err != nil {
		return err
	}
	var logEntry MedicineLogEntry
	if err := ctx.Bind(&logEntry); err != nil {
		return err
	}
	log.Println("Got a request to log medicine: ", logEntry)
	if logEntry.Time.IsZero() {
		logEntry.Time = time.Now()
	}
	if err := s.db.AddMedicineLog(user, logEntry); err != nil {
		return echo.ErrBadRequest
	} else {
		return ctx.JSON(http.StatusOK, nil)
	}
}

func (s *Server) GetMedicines(user userGUID) ([]MedicineType, error) {
	if medicines, err := s.db.GetMedicines(user); err != nil {
		return nil, err
	} else {
		var result []MedicineType
		for _, m := range medicines {
			result = append(result, m.MedicineType)
		}
		return result, nil
	}
}

func (s *Server) GetApiV1Medicines(ctx echo.Context) error {
	var user userGUID
	var err error
	if user, err = s.GetOrCreateUserGUIDCookie(ctx); err != nil {
		return err
	}
	if result, err := s.GetMedicines(user); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, result)
	}
}

func (s *Server) SetMedicines(user userGUID, medicines []MedicineType) error {
	for _, m := range medicines {
		if err := s.db.AddMedicine(
			MedicineTypeDB{
				MedicineType: m,
				User:         user,
			}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) PostApiV1Medicines(ctx echo.Context) error {
	var user userGUID
	var err error
	if user, err = s.GetOrCreateUserGUIDCookie(ctx); err != nil {
		return err
	}
	var medicines []MedicineType
	if err := ctx.Bind(&medicines); err != nil {
		return err
	}
	for _, m := range medicines {
		log.Println("Got medicine: ", m)
	}
	if err := s.SetMedicines(user, medicines); err != nil {
		log.Println("Failed to set medicines: ", err)
		return err
	}
	return s.GetApiV1Medicines(ctx)
	// return nil
}

func (s *Server) GetApiV1Logout(ctx echo.Context) error {
	// Set an invalid JWT cookie named "token" to log the user out
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(time.Hour * 24 * 365)
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	ctx.SetCookie(cookie)
	return ctx.JSON(http.StatusOK, nil)
}

func (s *Server) GetApiV1Settings(ctx echo.Context) error {
	var user userGUID
	var err error
	if user, err = s.GetOrCreateUserGUIDCookie(ctx); err != nil {
		return err
	}
	log.Println("Got a request for settings")
	if settings, err := s.db.GetSettings(user); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, settings)
	}
}

func (s *Server) PostApiV1Settings(ctx echo.Context) error {
	var user userGUID
	var err error
	if user, err = s.GetOrCreateUserGUIDCookie(ctx); err != nil {
		return err
	}
	var settings UserSettings
	if err := ctx.Bind(&settings); err != nil {
		return err
	}
	if err := s.db.UpdateSettings(user, settings); err != nil {
		return err
	}
	return nil
}

func (s *Server) GetApiV1DeleteUser(ctx echo.Context) error {
	var user userGUID
	var err error
	if user, err = s.GetOrCreateUserGUIDCookie(ctx); err != nil {
		return err
	}
	if err := s.db.DeleteUser(user); err != nil {
		return err
	}
	return nil
}
