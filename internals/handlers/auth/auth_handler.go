package auth

import (
	"database/sql"
	"encoding/json"
	"file-analyzer/internals/domain"
	"file-analyzer/internals/utils"
	"file-analyzer/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
	jwt.RegisteredClaims
}

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterDetails struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthHandler struct {
	l  *log.Logger
	db domain.DBRepo
}

var jwtKey = os.Getenv("JWT_SECRET_KEY")

func NewAuthHandler(l *log.Logger, db domain.DBRepo) *UserAuthHandler {
	return &UserAuthHandler{l, db}
}

func (cc *UserAuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userDetails LoginDetails
	err := DecodeJSON(r, &userDetails)
	if err != nil {
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if userDetails.Email == "" || userDetails.Password == "" {
		utils.FAIL(w, http.StatusBadRequest, "Email or Password is empty.")
		return
	}
	userId, err := cc.db.CheckUserExist(userDetails.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.FAIL(w, http.StatusNotFound, "User Not Found")
			return
		}
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// check for password
	user, err := cc.db.FindUserById(userId)
	if err != nil {
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(userDetails.Password))
	if err != nil {
		utils.FAIL(w, http.StatusUnauthorized, "Incorrect Credentials")
		return
	}
	expiry := time.Now().Add(5 * time.Minute)
	accessToken, err := GenerateJWT(userId, expiry)
	if err != nil {
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	refresh_token := ""
	SetCookie(r, w, "refresh_token", refresh_token, 7*time.Now().Day())
	utils.SUCCESS(w, "Login Successfully", accessToken)
}

func (cc *UserAuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userDetails RegisterDetails
	err := DecodeJSON(r, &userDetails)
	if err != nil {
		cc.l.Println(err)
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	user := models.User{
		Name:         userDetails.Name,
		Email:        userDetails.Email,
		PasswordHash: passwordHash,
	}
	err = cc.db.InsertUser(user)
	if err != nil {
		utils.FAIL(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	utils.SUCCESS(w, "Users Registered Successfully", nil)
}

func GenerateJWT(userId string, expiry time.Time) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}
	s := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := s.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func SetCookie(r *http.Request, w http.ResponseWriter, name string, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   r.TLS != nil,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	})
}

func DecodeJSON[T *LoginDetails | *RegisterDetails](r *http.Request, dst T) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
