package cart

import (
	"database/sql"
	"errors"
	"time"

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

func (s *CartServiceImpl) AddToCart(requestFormat CartItemRequestFormat, userID uuid.UUID) (cart Cart, err error)  {
	productID := requestFormat.ProductID
	quantity := requestFormat.Quantity

	productCheck , err := s.ProductService.ResolveProductByID(productID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if productCheck.Stock < int64(quantity) {
		err = failure.BadRequestFromString("Quantity cannot greater than stock")
		return
	}

	cart, err = s.CartRepository.ResolveCartByUserID(userID)
	if err != sql.ErrNoRows && err != nil {
		logger.ErrorWithStack(err)
		return
	} 
	if err == sql.ErrNoRows {
		cartID, errUUID := uuid.NewV4()
		if errUUID != nil {
			logger.ErrorWithStack(errUUID)
			return 
		}
	
		err = s.CartRepository.CreateCart(Cart{
			ID: cartID,
			UserID: userID,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		})
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
	} 

	existingItem, err := s.CartRepository.ResolveCartItemByProductID(cart.ID, productID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if existingItem == nil {
		if quantity <= 0 {
			err = failure.BadRequest(errors.New("quantity not valid"))
			return
		}
		
		err = s.CartRepository.CreateCartItem(CartItem{
			CartID: cart.ID,
			ProductID: productID,
			Quantity: quantity,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		})
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}


	} else {
		err = existingItem[0].Update(CartItem{ 
			CartID: cart.ID,
			ProductID: productID,
			Quantity: quantity,
		}, userID)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}

		err = s.CartRepository.UpdateCartItem(existingItem[0])
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
	}

	items, err := s.CartRepository.ResolveCartItemsJoinProduct(cart.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	cart.AttachItems(items)
	cart.Update(userID)
	err = s.CartRepository.UpdateCart(cart)
	if err != nil {
		logger.ErrorWithStack(err)
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
	exists, err := s.CartRepository.ExistsByID(cartID)
	if err != nil {
		return
	}
	if !exists {
		err = failure.BadRequestFromString("cart not found")
		logger.ErrorWithStack(err)
		return
	}

	isHaveAccess, err := s.checkCartOwner(cartID, userID, role)
	if err != nil {
		return
	}
	if !isHaveAccess {
		err = failure.Unauthorized("untauthorized")
		return
	}

	orderID, err := uuid.NewV4()
	if err != nil {
		return
	}
	newOrder = order.Order{
		ID: orderID,
		UserID: userID,
		Address: requestFormat.Address,
		Status: "pending",
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	var orderItems = []order.OrderItem{}
	for _, item := range requestFormat.ProductIDs {
		cartItem, err := s.CartRepository.ResolveCartItemJoinProduct(cartID, item)
		if err != nil {
			logger.ErrorWithStack(err)
			return order.Order{}, err
		}

		if cartItem.Quantity > cartItem.Stock {
			log.Error().Msg("out of stock")
			err = errors.New("out of stock")
			return order.Order{}, err
		}
		
		orderItem := order.OrderItem{
			OrderID: orderID,
			ProductID: cartItem.ProductID,
			Quantity: cartItem.Quantity,
			UnitPrice: cartItem.UnitPrice,
			CreatedAt: time.Now(),
			CreatedBy: userID,	
		}
		orderItems = append(orderItems, orderItem)
	}

	newOrder.AttachItems(orderItems)
	newOrder.Recalculate()

	// create order
	err = s.OrderService.CreateOrder(newOrder)
	if err != nil {
		return
	}
	// create order items and update stock 
	for _, orderItem := range newOrder.Items {
		err = s.OrderService.CreateOrderItem(orderItem)
		if err != nil {
			return
		}
		// hard delete cart_item 
		err = s.CartRepository.DeleteCartItem(cartID, orderItem.ProductID)
		if err != nil {
			return
		}
	}	

	return

}

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