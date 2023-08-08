package handlers

import (
	"net/http"
	"strconv"

	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductService product.ProductService
	AuthMiddleware *middleware.Authentication
}

func ProvideProductHandler(productService product.ProductService, authMiddleware *middleware.Authentication) ProductHandler {
	return ProductHandler{
		ProductService: productService,
		AuthMiddleware: authMiddleware,
	}
}

func (h *ProductHandler) Router(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", h.ResolveAllProducts)
		})

		// r.Group(func(r chi.Router) {
		// 	r.Use(h.AuthMiddleware.Password)
		// 	r.Post("/foo", h.CreateFoo)
		// 	r.Delete("/foo/{id}", h.SoftDeleteFoo)
		// 	r.Put("/foo/{id}", h.UpdateFoo)
		// })

	})
}

func (h *ProductHandler) ResolveAllProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Must have page queryparam")
		return	
	}
	limitStr := r.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Must have limit queryparam")
		return	
	}

	products, err := h.ProductService.ResolveAllProducts(pageInt, limitInt)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, products)
}



