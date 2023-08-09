package cart

import (
	"database/sql"
	"time"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/product"
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
	ProductService product.ProductService
}

func ProvideCartServiceImpl(cartRepository CartRepository, conf *configs.Config, productService product.ProductService) *CartServiceImpl  {
	s := new(CartServiceImpl)
	s.CartRepository = cartRepository
	s.ProductService = productService
	s.Config = conf

	return s
}

func (s *CartServiceImpl) AddToCart(requestFormat CartItemRequestFormat, userID uuid.UUID) (cart Cart, err error)  {
	productID := requestFormat.ProductID
	quantity := requestFormat.Quantity


	_ , err = s.ProductService.ResolveProductByID(productID)
	if err != nil {
		logger.ErrorWithStack(err)
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

	for i := 0; i < len(items) ; i++ {
		product, errProduct := s.ProductService.ResolveProductByID(items[i].ProductID)
		if errProduct != nil {
			logger.ErrorWithStack(errProduct)
			return
		}
		items[i].UnitPrice = product.Price
		items[i].Recalculate()
	}

	cart.AttachItems(items)
	cart.Update(userID)

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

	for i := 0; i < len(items) ; i++ {
		product, errProduct := s.ProductService.ResolveProductByID(items[i].ProductID)
		if errProduct != nil {
			logger.ErrorWithStack(errProduct)
			return
		}
		items[i].UnitPrice = product.Price
		items[i].Recalculate()
	}

	cart.AttachItems(items)
	cart.Recalculate()
	return

}