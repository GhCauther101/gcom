package user

import (
	"fmt"
	"gcom/config"
	"gcom/service/auth"
	"gcom/types"
	"gcom/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	//parse JSON payload
	var loginPayload types.LoginUserPayload

	if err := utils.ReadJSON(r, &loginPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[api register json] failed to parse request body\n", err))
		return
	}

	if err := utils.Validate.Struct(loginPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload &v", err))
		return
	}

	u, err := h.store.GetUserByEmail(loginPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found, invalid email or password"))
		return
	}

	if !auth.ComparePassword(u.Password, loginPayload.Password) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}
	
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	//parse JSON payload
	var registerPayload types.RegisterUserPayload

	if err := utils.ReadJSON(r, &registerPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[api register json] failed to parse request body\n", err))
		return
	}

	if err := utils.Validate.Struct(registerPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload &v", err))
		return
	}

	_, err := h.store.GetUserByEmail(registerPayload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[api] user with email already exists %s", registerPayload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(registerPayload.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	err = h.store.CreateUser(types.User{
		FirstName: registerPayload.FirstName,
		LastName:  registerPayload.LastName,
		Email:     registerPayload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
