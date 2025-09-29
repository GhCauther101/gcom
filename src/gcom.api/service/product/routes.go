package product

import (
	"gcom/types"
	"gcom/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.HandleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", h.HandleCreateProduct).Methods(http.MethodPost)
}

func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	
}

func (h *Handler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	productStore, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, productStore)
}
