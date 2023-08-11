package handlers

import (
	"net/http"
	"strconv"

	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type OrderHandler struct {
	OrderService   order.OrderService
	AuthMiddleware *middleware.Authentication
}

func ProvideOrderHandler(orderService order.OrderService, authMiddleware *middleware.Authentication) OrderHandler {
	return OrderHandler{
		OrderService: orderService,
		AuthMiddleware: authMiddleware,
	}
}


func (h *OrderHandler) Router(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware.ValidateJWT)
			r.Get("/", h.ResolveAllOrder)
		})

	})
}

func (h *OrderHandler) ResolveAllOrder(w http.ResponseWriter, r *http.Request)  {
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
	
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(shared.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
	}

	id, err := uuid.FromString(claims.UserId)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	orders, err := h.OrderService.ResolveAllOrder(id, claims.Role, pageInt, limitInt)


	response.WithJSON(w, http.StatusOK, orders)
}


