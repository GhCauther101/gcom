package cart

import (
	"fmt"
	"gcom/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}

		productIds[i] = item.ProductID	
	}

	return productIds, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	if err := checkIfCartIsStock(items, productMap); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(items, productMap) 	

	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProduct(product)
	}

	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID: userID,
		Total: int(totalPrice),
		Status: "pending",
		Address: "rnd address",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID: orderID,
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			Price: productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) <= 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please update cart containments", item.ProductID)
		}

		if product.Quantity < item.ProductID {
			return fmt.Errorf("product %s is not available", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var total float64
	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(product.Quantity)
	}

	return total
}