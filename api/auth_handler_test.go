package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/vishnupal/hotel-reservation/db"
	"github.com/vishnupal/hotel-reservation/types"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "james@foo.com",
		FirstName: "james",
		LastName:  "foo",
		Password:  "supersecurepassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	_ = insertTestUser(t, tdb.UserStore)

	app := fiber.New()

	authHandler := NewAuthHandler(tdb.UserStore)

	app.Post("/auth", authHandler.HandleAunthenticate)

	parmas := AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepasswordnotcorrect",
	}

	b, _ := json.Marshal(parmas)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 200 but  got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}
	if genResp.Msg != "Invalid credentials" {
		t.Fatalf(
			"expected gen response Message to be <Invalid credentials> but got %s",
			genResp.Msg,
		)
	}
}

func TestAuthenticationSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)

	app := fiber.New()

	authHandler := NewAuthHandler(tdb.UserStore)

	app.Post("/auth", authHandler.HandleAunthenticate)

	parmas := AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepassword",
	}

	b, _ := json.Marshal(parmas)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 400 but  got %d", resp.StatusCode)
	}

	var authResp AuthReponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}

	// set the expected Password to an Empty string, beacuse we do Not return nay
	// JSON response .
	insertedUser.EncyptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResp.User)
		t.Fatalf("expected the user to be the inserted user")
	}
}
