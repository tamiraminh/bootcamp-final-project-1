package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
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

		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware.ValidateJWT)
			r.Use(h.AuthMiddleware.RoleAdminCheck)
			r.Post("/", h.CreateProduct)
		})

	})
}

// CreateFoo creates a new Product.
// @Summary Create a new Product.
// @Description This endpoint creates a new Product.
// @Tags v1/Products
// @Security JWTToken
// @Param foo body product.ProductRequestFormat true "The Product to be created."
// @Produce json
// @Success 201 {object} response.Base{data=product.ProductResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat product.ProductRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(shared.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
	}

	id, err := uuid.FromString(claims.UserId)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	product, err := h.ProductService.Create(requestFormat, id )
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, product)
}


// @Summary Resolve All Products
// @Description This endpoint resolves All products with pagination page and limit.
// @Tags v1/Products
// @Param page query int true "must greater or equeal to zero"
// @Param limit query int true "must greater than zero"
// @Produce json
// @Success 200 {object} response.Base{data=product.ProductResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/products [get]
func (h *ProductHandler) ResolveAllProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Must have page queryparam")
		return	
	}
	if pageInt < 0 {
		response.WithMessage(w, http.StatusBadRequest, "page must be equal or greater to zero")
		return
	}
	limitStr := r.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		response.WithMessage(w, http.StatusBadRequest, "Must have limit queryparam")
		return	
	}
	if limitInt <= 0 {
		response.WithMessage(w, http.StatusBadRequest, "limit must be greater to zero")
		return
	}
	products, err := h.ProductService.ResolveAllProducts(pageInt, limitInt)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, products)
}



