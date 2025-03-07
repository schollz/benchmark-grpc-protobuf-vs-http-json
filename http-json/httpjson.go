package httpjson

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
)

// Start entrypoint
func Start() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/", func(c *gin.Context) {
		var user User
		c.Bind(&user)

		validationErr := validate(user)
		if validationErr != nil {
			c.JSON(200, Response{
				Code:    500,
				Message: validationErr.Error(),
			})
			return
		}

		user.ID = "1000000"
		c.JSON(200, Response{
			Code:    200,
			Message: "OK",
			User:    &user,
		})
	})
	r.Run(":60001")
	// http.HandleFunc("/", CreateUser)
	// log.Println(http.ListenAndServe(":60001", nil))
}

// User type
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Response type
type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	User    *User  `json:"user"`
}

// CreateUser handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	decoder.Decode(&user)
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	validationErr := validate(user)
	if validationErr != nil {
		json.NewEncoder(w).Encode(Response{
			Code:    500,
			Message: validationErr.Error(),
		})
		return
	}

	user.ID = "1000000"
	json.NewEncoder(w).Encode(Response{
		Code:    200,
		Message: "OK",
		User:    &user,
	})
}

func validate(in User) error {
	_, err := mail.ParseAddress(in.Email)
	if err != nil {
		return err
	}

	if len(in.Name) < 4 {
		return errors.New("Name is too short")
	}

	if len(in.Password) < 4 {
		return errors.New("Password is too weak")
	}

	return nil
}
