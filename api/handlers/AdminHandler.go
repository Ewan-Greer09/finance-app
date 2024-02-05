package handlers

import (
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"

	"github.com/Ewan-Greer09/finance-app/api/database"
	"github.com/Ewan-Greer09/finance-app/api/models"
)

type AdminHandler struct {
	Logger *slog.Logger
	DB     database.Database
	FS     embed.FS
}

func NewAdminHandler(log *slog.Logger, db database.Database, fs embed.FS) *AdminHandler {
	return &AdminHandler{
		Logger: log,
		DB:     db,
		FS:     fs,
	}
}

func (a *AdminHandler) Routes(r chi.Router) {
	r.Post("/login", a.Login)
	r.Get("/logout", a.Logout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(a.IsAdmin) // JWT middleware
		r.Get("/user", a.GetUser)
		r.Post("/user", a.CreateUser)
	})
}

func (a *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if username == "" {
		render.HTML(w, r, "<h1>Username is required</h1>")
		render.Status(r, http.StatusBadRequest)
		return
	}

	user, err := a.DB.GetUser(username)
	if err != nil {
		render.HTML(w, r, "<h1>User not found</h1>")
		render.Status(r, http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFS(a.FS, "web/components/user.html")
	if err != nil {
		a.Logger.Error(parseTemplateError, "error", err)
		http.Error(w, parseTemplateError, http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, user)
	if err != nil {
		a.Logger.Error(executeTemplateError, "error", err)
		http.Error(w, executeTemplateError, http.StatusInternalServerError)
	}
}

func (a *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	isAdmin := r.FormValue("isAdmin")
	if isAdmin == "true" || isAdmin == "on" {
		user.IsAdmin = true
	}

	if user.Username == "" || user.Password == "" {
		render.HTML(w, r, "<h1>Username and Password are required</h1>")
		render.Status(r, http.StatusBadRequest)
		return
	}

	err := a.DB.CreateUser(user)
	if err != nil {
		render.HTML(w, r, "<h1>Failed to create user</h1>")
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.HTML(w, r, "<h1>User created</h1>")
}

func (a *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		render.HTML(w, r, "<h1>Username and Password are required</h1>")
		render.Status(r, http.StatusBadRequest)
		return
	}

	user, err := a.DB.GetUser(username)
	if err != nil {
		render.HTML(w, r, "<h1>Invalid Username or Password</h1>")
		render.Status(r, http.StatusUnauthorized)
		return
	}
	if user.Password == "" || user.Username == "" {
		render.HTML(w, r, "<h1>Invalid Username or Password</h1>")
		render.Status(r, http.StatusUnauthorized)

		return
	}

	if password != user.Password {
		render.HTML(w, r, "<h1>Invalid Username or Password</h1>")
		render.Status(r, http.StatusUnauthorized)

		return
	}

	cookie := r.Cookies()
	for _, c := range cookie {
		if c.Name == "access-token" {
			render.HTML(w, r, "<h1>Already logged in</h1>")
			render.Status(r, http.StatusAccepted)
			return
		}
	}

	accessToken, exp, err := generateAccessToken(&user)
	if err != nil {
		render.HTML(w, r, "<h1>Internal Server Error</h1>")
		render.Status(r, http.StatusInternalServerError)
		return
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, w)
	setUserCookie(&user, exp, w)

	render.HTML(w, r, "<h1>Logged in</h1>")
}

func (a *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logoutTokenCookie(accessTokenCookieName, "", w)
	render.HTML(w, r, "<h1>Logged out</h1>")
}

// middleware to check if user is admin
func (a *AdminHandler) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := r.Cookies()
		for _, c := range cookie {
			if c.Name == "access-token" {
				// get token from cookie
				token, err := jwt.Parse(c.Value, func(token *jwt.Token) (interface{}, error) {
					// Don't forget to validate the alg is what you expect:
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("Unexpected signing method")
					}
					return []byte(GetJWTSecret()), nil
				})
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok || !token.Valid {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				// check if user is admin
				user, err := a.DB.GetUser(claims["name"].(string))
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				if user.IsAdmin != true {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				next.ServeHTTP(w, r)
			}
		}
	})
}
