package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(user User) error
}

type OrderStore interface {
	CreateOrder(order Order) (int, error)
	CreateOrderItem(orderItem OrderItem) (error)
}

type Order struct {
	ID int `json:"id"`
	UserID int `json:"userID"`
	Total int `json:"total"`
	Status string `json:"status"`
	Address string `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID int `json:"id"`
	OrderID int `json:"orderID"`
	ProductID int `json:"productID"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProductStore interface {
	GetProducts() ([]Product, error)
	GetProductsByIDs(ps []int) ([]Product, error)
	UpdateProduct(product Product) (error)
}

type Product struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	Price float64 `json:"price"`
	Quantity int `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=100"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CartItem struct {
	ProductID int `json:"productID"`
	Quantity int `json:"quantity"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}