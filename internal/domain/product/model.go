package product

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Product struct {
	ID            uuid.UUID   `db:"id" validate:"required"`
	Name          string      `db:"name" validate:"required"`
	Description	  string      `db:"description" validate:"required"`
	Category	  string      `db:"category" validate:"required"`
	Brand		  string      `db:"brand" validate:"required"`
	Stock		  int64       `db:"stock" validate:"required"`
	Price		  float64     `db:"price" validate:"required"`
	CreatedAt     time.Time   `db:"created_at" validate:"required"`
	CreatedBy     uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt     null.Time   `db:"updated_at"`
	UpdatedBy     nuuid.NUUID `db:"updated_by"`
	DeletedAt     null.Time   `db:"deleted_at"`
	DeletedBy     nuuid.NUUID `db:"deleted_by"`
}

func (p *Product) IsDeleted() (deleted bool) {
	return p.DeletedAt.Valid && p.DeletedBy.Valid
}

func (p Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponseFormat())
}


func (p Product) NewFromRequestFormat(req ProductRequestFormat, userID uuid.UUID) (newProduct Product, err error) {
	productID, _ := uuid.NewV4()
	newProduct = Product{
		ID:          productID,
		Name:        req.Name,
		Description: req.Description,
		Category: req.Category,
		Brand: req.Brand,
		Stock: req.Stock,
		Price: req.Price,
		CreatedAt:     time.Now(),
		CreatedBy:   userID,
	}

	err = newProduct.Validate()

	return
}


func (p Product) ToResponseFormat() ProductResponseFormat {
	resp := ProductResponseFormat{
		ID:            	p.ID,
		Name:          	p.Name,
		Description: 	p.Description,
		Category: 		p.Category,
		Brand: 			p.Brand,
		Stock: 			p.Stock,
		Price: 			p.Price,
		CreatedAt:      p.CreatedAt,
		CreatedBy:     	p.CreatedBy,
		UpdatedAt:      p.UpdatedAt,
		UpdatedBy:     	p.UpdatedBy.Ptr(),
		DeletedAt:      p.DeletedAt,
		DeletedBy:     	p.DeletedBy.Ptr(),
	}


	return resp
}


func (p *Product) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(p)
}


type ProductRequestFormat struct {
	Name 			string   `json:"name" validate:"required"`
	Description 	string	 `json:"description" validate:"required"`
	Category		string	 `json:"category" validate:"required"`
	Brand 			string 	 `json:"brand" validate:"required"`
	Stock			int64	 `json:"stock" validate:"required"`
	Price			float64  `json:"price" validate:"required"`
}


type ProductResponseFormat struct {
	ID            uuid.UUID   `db:"id" validate:"required"`
	Name 			string   `json:"name" validate:"required"`
	Description 	string	 `json:"description" validate:"required"`
	Category		string	 `json:"category" validate:"required"`
	Brand 			string 	 `json:"brand" validate:"required"`
	Stock			int64	 `json:"stock" validate:"required"`
	Price			float64   `json:"price" validate:"required"`
	CreatedAt     	time.Time   `json:"created_at" validate:"required"`
	CreatedBy     	uuid.UUID   `json:"created_by" validate:"required"`
	UpdatedAt     	null.Time   `json:"updated_at"`
	UpdatedBy     	*uuid.UUID `json:"updated_by"`
	DeletedAt     	null.Time   `json:"deleted_at"`
	DeletedBy     	*uuid.UUID `json:"deleted_by"`
}



