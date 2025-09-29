package api

import (
	"database/sql"
	"gcom/service/cart"
	"gcom/service/order"
	"gcom/service/product"
	"gcom/service/user"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	productStore := product.NewStore(s.db)	
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subRouter)
	
	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subRouter)
	
	log.Println("[api-server] api server listenning on", s.addr)
	return http.ListenAndServe(s.addr, router)
}