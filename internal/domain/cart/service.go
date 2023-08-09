package cart

import (
	"database/sql"
	"time"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
)

type CartService interface {
	AddToCart(requestFormat CartItemRequestFormat, userID uuid.UUID) (cart Cart, err error)
	ResolveCartByUserID(userID uuid.UUID) (cart Cart, err error)
}

type CartServiceImpl struct {
	CartRepository CartRepository
	Config         *configs.Config
}

func ProvideCartServiceImpl(cartRepository CartRepository, conf *configs.Config) *CartServiceImpl  {
	s := new(CartServiceImpl)
	s.CartRepository = cartRepository
	s.Config = conf

	return s
}

func (s *CartServiceImpl) AddToCart(requestFormat CartItemRequestFormat, userID uuid.UUID) (cart Cart, err error)  {
	productID := requestFormat.ProductID
	quantity := requestFormat.Quantity

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
	
		cart = Cart{
			ID: cartID,
			UserID: userID,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		}
	
		err = s.CartRepository.CreateCart(cart)
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
		cartItem := CartItem{
			CartID: cart.ID,
			ProductID: productID,
			Quantity: quantity,
			CreatedAt: time.Now(),
			CreatedBy: userID,
		}

		err = s.CartRepository.CreateCartItem(cartItem)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}


	} else {
		cart.Update(userID)
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

	items, err := s.CartRepository.ResolveCartItemsByCartID(cart.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	cart.AttachItems(items)
	return

}


func (s *CartServiceImpl) ResolveCartByUserID(userID uuid.UUID) (cart Cart, err error)  {
	cart, err = s.CartRepository.ResolveCartByUserID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	} 


	items, err := s.CartRepository.ResolveCartItemsByCartID(cart.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	cart.AttachItems(items)
	return

}