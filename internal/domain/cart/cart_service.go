package cart

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type CartService interface {
	AddToCart(requestFormat CartItemRequestFormat, userID uuid.UUID) (cart Cart, err error)
	ResolveCartByUserID(userID uuid.UUID) (cart Cart, err error)
	Checkout(requestFormat CheckoutRequestFormat, userID uuid.UUID, cartID uuid.UUID, role string) (order order.Order, err error)
}

type CartServiceImpl struct {
	CartRepository CartRepository
	Config         *configs.Config
	ProductService product.ProductService
	OrderService   order.OrderService
}

func ProvideCartServiceImpl(cartRepository CartRepository, conf *configs.Config, productService product.ProductService, orderService order.OrderService) *CartServiceImpl  {
	s := new(CartServiceImpl)
	s.CartRepository = cartRepository
	s.ProductService = productService
	s.OrderService = orderService
	s.Config = conf

	return s
}

func (s *CartServiceImpl) AddToCart(req CartItemRequestFormat, userID uuid.UUID) (cart Cart, err error) {
	// Check product availability and stock
	product, err := s.ProductService.ResolveProductByID(req.ProductID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if product.Stock < int64(req.Quantity) {
		err = failure.BadRequestFromString("Quantity cannot be greater than stock")
		return
	}

	// Resolve or create cart for the user
	cart, err = s.resolveOrCreateCart(userID)
	if err != nil {
		return
	}

	// Handle cart item addition or update
	err = s.handleCartItem(cart, req , userID)
	if err != nil {
		return
	}

	// Update cart and return
	err = s.updateCart(&cart, userID)
	if err != nil {
		return
	}

	return
}




func (s *CartServiceImpl) ResolveCartByUserID(userID uuid.UUID) (cart Cart, err error)  {
	cart, err = s.CartRepository.ResolveCartByUserID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	} 


	items, err := s.CartRepository.ResolveCartItemsJoinProduct(cart.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	cart.AttachItems(items)
	cart.Recalculate()
	return
}

func (s *CartServiceImpl) Checkout(requestFormat CheckoutRequestFormat, userID uuid.UUID, cartID uuid.UUID, role string) (newOrder order.Order, err error) {
	// Check if cart exists
	if exists, err := s.CartRepository.ExistsByID(cartID); err != nil {
		return newOrder, err
	} else if !exists {
		err = failure.BadRequestFromString("cart not found")
		logger.ErrorWithStack(err)
		return newOrder, err
	}

	// Check cart owner access
	if isHaveAccess, err := s.checkCartOwner(cartID, userID, role); err != nil {
		return newOrder, err
	} else if !isHaveAccess {
		err = failure.Unauthorized("unauthorized")
		return newOrder, err
	}

	if err != nil {
		return newOrder, err
	}

	newOrder, err = order.Order{}.NewOrder(userID, requestFormat.Address)
	if err != nil {
		return newOrder, err
	}

	orderItems, err := s.createOrderItems(cartID, userID, newOrder.ID, requestFormat.ProductIDs)
	if err != nil {
		return newOrder, err
	}
	fmt.Println(orderItems)
	newOrder.AttachItems(orderItems)
	newOrder.Recalculate()

	if err = s.createOrderAndHandleCart(cartID, newOrder, orderItems); err != nil {
		return newOrder, err
	}

	return newOrder, nil
}


// internal function

func (s *CartServiceImpl) checkCartOwner(cartID uuid.UUID, userID uuid.UUID, role string) (isHaveAccess bool, err error) {
	cart, err := s.CartRepository.ResolveCartByID(cartID)
	if err != nil{
		return
	}
	if userID != cart.UserID && role != "admin" {
		isHaveAccess = false
		return
	}
	return true, err
}

func (s *CartServiceImpl) resolveOrCreateCart(userID uuid.UUID) (cart Cart, err error) {
	cart, err = s.CartRepository.ResolveCartByUserID(userID)
	if err == sql.ErrNoRows {
		cart, err = s.CartRepository.CreateCartByUserID(userID)
	}
	return
}

func (s *CartServiceImpl) handleCartItem(cart Cart, req CartItemRequestFormat, userID uuid.UUID) (err error) {
	existingItem, err := s.CartRepository.ResolveCartItemByProductID(cart.ID, req.ProductID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if len(existingItem) == 0 {
		if req.Quantity <= 0 {
			err = failure.BadRequest(errors.New("quantity not valid"))
			return
		}

		newCartItem := CartItem{}.newCartItem(cart.ID, req, userID)

		err = s.CartRepository.CreateCartItem(newCartItem)
		if err != nil {
			return
		}

	} else {
		updateCartItem := CartItem{}.newCartItem(cart.ID, req, userID)
		err = existingItem[0].Update(updateCartItem, cart.UserID)
		if err == nil {
			err = s.CartRepository.UpdateCartItem(existingItem[0])
		}
	}
	return
}


func (s *CartServiceImpl) updateCart(cart *Cart, userID uuid.UUID) (err error) {
	items, err := s.CartRepository.ResolveCartItemsJoinProduct(cart.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	cart.AttachItems(items)
	cart.Update(userID)

	err = s.CartRepository.UpdateCart(*cart)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (s *CartServiceImpl) createOrderItems(cartID uuid.UUID, userID uuid.UUID, orderID uuid.UUID, productIDs []uuid.UUID) ([]order.OrderItem, error) {
	var orderItems []order.OrderItem


	for _, item := range productIDs {
		cartItem, err := s.CartRepository.ResolveCartItemJoinProduct(cartID, item)
		if err != nil {
			logger.ErrorWithStack(err)
			return nil, err
		}

		if cartItem.Quantity > cartItem.Stock {
			err = errors.New("out of stock")
			log.Error().Msg("out of stock")
			return nil, err
		}

		orderItem := order.OrderItem{}.NewOrderItem(orderID, userID, cartItem.ProductID, cartItem.Quantity, cartItem.UnitPrice)
		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}


func (s *CartServiceImpl) createOrderAndHandleCart(cartID uuid.UUID, newOrder order.Order, orderItems []order.OrderItem) (err error) {
	// Create order
	if err := s.OrderService.CreateOrder(newOrder); err != nil {
		return err
	}
	// Create order items and update stock
	for _, orderItem := range orderItems {
		if err := s.OrderService.CreateOrderItem(orderItem); err != nil {
			return err
		}

		if err := s.CartRepository.DeleteCartItem(cartID, orderItem.ProductID); err != nil {
			return err
		}
	}

	return
}

