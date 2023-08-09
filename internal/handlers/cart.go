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
		})


	})
}


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
