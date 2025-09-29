package test

import (
	"bytes"
	"encoding/json"
	"fmt"

	user "gcom/service/user"
	"gcom/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := user.NewHandler(userStore)
	
	t.Run("should fail if user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload {
			FirstName: "user",
			LastName: "123",
			Email: "asd",
			Password: "asd",
		}
		marshalled, _ := json.Marshal(payload)
		
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(r, req)
		
		if r.Code != http.StatusBadRequest {
			t.Errorf("[test] expected status code %d, got %d", http.StatusBadRequest, r.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload {
			FirstName: "user",
			LastName: "123",
			Email: "asdv@gmail.com",
			Password: "asd",
		}
		marshalled, _ := json.Marshal(payload)
		
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(r, req)
		
		if r.Code != http.StatusCreated {
			t.Errorf("[test] expected status code %d, got %d", http.StatusCreated, r.Code)
		}
	})
}

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}