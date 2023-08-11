package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type CartHandler struct {
	CartService cart.CartService
	AuthMiddleware *middleware.Authentication
}

func ProvideCartHandler(cartService cart.CartService, authmiddleware *middleware.Authentication) CartHandler {
	return CartHandler{
		CartService: cartService,
		AuthMiddleware: authmiddleware,
	}
}


func (h *CartHandler) Router(r chi.Router) {
	r.Route("/carts", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware.ValidateJWT)
			r.Post("/", h.AddToCart)
			r.Get("/", h.ResolveCartByUserID)
			r.Post("/{cart_id}/checkout", h.Checkout)
		})


	})
}

// @Summary Add product to cart.
// @Description This endpoint is to add product to cart.
// @Tags v1/Cart
// @Security JWTToken
// @Param Cart body cart.CartItemRequestFormat true "make add to cart request"
// @Produce json
// @Success 201 {object} response.Base{data=cart.CartResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/carts [post]
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat cart.CartItemRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
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

	cart,err := h.CartService.AddToCart(requestFormat, id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, cart)
}

// @Summary Resolve Cart by user ID
// @Description This endpoint resolves a Cart by its user ID.
// @Tags v1/Cart
// @Security JWTToken
// @Produce json
// @Success 200 {object} response.Base{data=cart.CartResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/carts [get]
func (h *CartHandler) ResolveCartByUserID(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(shared.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
	}

	id, err := uuid.FromString(claims.UserId)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}


	cart,err := h.CartService.ResolveCartByUserID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, cart)
	
}


// @Summary checkout selected Product.
// @Description This endpoint checkout selected product in cart.
// @Tags v1/Cart
// @Security JWTToken
// @Param cart_id path string true "cartID"
// @Param Order body cart.CheckoutRequestFormat true "make order request"
// @Produce json
// @Success 201 {object} response.Base{data=order.OrderResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/carts/{cart_id}/checkout [post]
func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request)  {
	cartString := chi.URLParam(r, "cart_id")
	cartID, err := uuid.FromString(cartString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestFormat cart.CheckoutRequestFormat
	err = decoder.Decode(&requestFormat)
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

	order, err := h.CartService.Checkout(requestFormat, id, cartID, claims.Role )
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, order)
}
