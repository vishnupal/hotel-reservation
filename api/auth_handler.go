package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vishnupal/hotel-reservation/db"
	"github.com/vishnupal/hotel-reservation/types"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthReponse struct {
	User  *types.User `json:"email"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invlidCredentails(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(
		genericResp{
			Type: "error",
			Msg:  "Invalid credentials",
		})
}

// A Handler should only do:
// - serialization of the incoming request (JSON). bytes data covert into Object
// - do some data fetching from db
// - call some business logic
// - return data back to user

func (h *AuthHandler) HandleAunthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invlidCredentails(c)
		}

		return err
	}
	if !types.IsValidPassword(user.EncyptedPassword, params.Password) {
		return invlidCredentails(c)
	}
	token := createTokenFromUser(user)
	resp := AuthReponse{
		User:  user,
		Token: token,
	}
	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenString
}
