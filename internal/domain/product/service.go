package product

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type ProductService interface {
	Create(requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error)
	ResolveAllProducts(page int, limit int) (product []Product, err error)
}

type ProductServiceImpl struct {
	ProductRepository 	ProductRepository
	Config				*configs.Config
}

func ProvideProductServiceImpl(productRepository ProductRepository, config *configs.Config) *ProductServiceImpl {
	s := new(ProductServiceImpl)
	s.ProductRepository = productRepository
	s.Config = config

	return s
}

func (s *ProductServiceImpl) Create(requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error) {
	product, err = product.NewFromRequestFormat(requestFormat, userID)
	if err != nil {
		return
	}

	if err != nil {
		return product, failure.BadRequest(err)
	}

	err = s.ProductRepository.Create(product)

	if err != nil {
		return
	}


	return
}


func (s *ProductServiceImpl) ResolveAllProducts(page int, limit int) (products []Product, err error) {
	products, err = s.ProductRepository.ResolveAllProducts(page, limit)
	if err != nil {
		return
	}

	return
}




