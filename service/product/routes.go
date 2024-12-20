package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kidusshun/ecom_bot/service/user"
	"github.com/kidusshun/ecom_bot/utils"
)

type Handler struct {
	store     ProductStore
	userStore user.UserStore
}

func NewHandler(store ProductStore, userStore user.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/products", h.handleGetProducts)
	router.Get("/products/{productID}", h.handleGetProductByID)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, 200, products)
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	id, err := uuid.Parse(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := h.store.GetProductByID(id)
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}

	utils.WriteJSON(w, 200, product)

}
