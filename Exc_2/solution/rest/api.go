package rest

import (
	"net/http"
	"ordersystem/repository"

	"encoding/json"
	"ordersystem/model"
	"time"

	"github.com/go-chi/render"
)

// GetMenu 			godoc
// @tags 			Menu
// @Description 	Returns the menu of all drinks
// @Produce  		json
// @Success 		200 {array} model.Drink
// @Router 			/api/menu [get]
func GetMenu(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo
		// get slice from db
		// render.Status(r, http.StatusOK)
		// render.JSON(w, r, <your-slice>)
		getDrinks := db.GetDrinks()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, getDrinks)
	}
}

// todo create GetOrders /api/order/all

// GetMenu 			godoc
// @tags 			AllOrders
// @Description 	Return all the placed orders
// @Produce  		json
// @Success 		200 {array} model.Order
// @Router 			/api/order/all [get]
func GetOrders(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getOrders := db.GetOrders()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, getOrders)
	}
}

// todo create GetOrdersTotal /api/order/total

// GetMenu 			godoc
// @tags 			GetTotalOrders
// @Description 	Return how many times each drink was purchased
// @Produce  		json
// @Success 		200 {array} model.Order
// @Router 			/api/order/totalled [get]
func GetOrdersTotal(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gettotalOrders := db.GetTotalledOrders()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, gettotalOrders)
	}
}

// PostOrder 		godoc
// @tags 			Order
// @Description 	Adds an order to the db
// @Accept 			json
// @Param 			b body model.Order true "Order" example({"drink_id":1,"amount":2}) //Google is my best friend
// @Produce  		json
// @Success 		200
// @Failure     	400
// @Router 			/api/order [post]
func PostOrder(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo
		// declare empty order struct
		// err := json.NewDecoder(r.Body).Decode(&<your-order-struct>)
		// handle error and render Status 400
		// add to db
		var order model.Order

		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Bad Request")
			return
		}
		order.CreatedAt = time.Now()

		db.AddOrder(&order)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok")
	}
}
