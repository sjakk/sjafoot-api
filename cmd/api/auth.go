package main

import (
	"encoding/json"
	"net/http"
	_ "github.com/sjakk/sjafoot/internal/data"
	"time"

	"github.com/sjakk/sjafoot/internal/data"
    	"github.com/lib/pq"

	"github.com/golang-jwt/jwt/v5"
)

// loginHandler handles user authentication and token generation.
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil || !match {
		app.errorResponse(w, r, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]string{
		"token": tokenString,
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}



func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	

	count, err := app.models.Users.Count()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}


	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: true,
	}
		

	if count == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}


	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := data.NewValidator()

	if v.Check(input.Name != "", "name", "must be provided"); !v.Valid() {
		// Further checks can be added here
	}
	if v.Check(input.Email != "", "email", "must be provided"); !v.Valid() {
	} else if v.Check(data.Matches(input.Email, data.EmailRX), "email", "must be a valid email address"); !v.Valid() {
	}
	if v.Check(input.Password != "", "password", "must be provided"); !v.Valid() {
	} else if v.Check(len(input.Password) >= 8, "password", "must be at least 8 bytes long"); !v.Valid() {
	}


	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			v.AddError("email", "an account with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]string{
		"token": tokenString,
	}
	app.writeJSON(w, http.StatusCreated, response, nil)
}
