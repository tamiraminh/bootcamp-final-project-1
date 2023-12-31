package cart

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Cart struct {
	ID 			uuid.UUID   `db:"id" validate:"required"`
	UserID 		uuid.UUID 	`db:"user_id" validate:"required"`
	CreatedAt	time.Time   `db:"created_at" validate:"required"`
	CreatedBy	uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt	null.Time   `db:"updated_at"`
	UpdatedBy	nuuid.NUUID `db:"updated_by"`
	DeletedAt	null.Time   `db:"deleted_at"`
	DeletedBy	nuuid.NUUID `db:"deleted_by"`
	TotalPrice	float64		`db:"-"`
	Items		[]CartItem  `db:"-"`
}


func (c *Cart) AttachItems(items []CartItem) Cart {
	for _, item := range items {
		if item.CartID == c.ID {
			c.Items = append(c.Items, item)
		}
	}

	return *c
}

func (c *Cart) Recalculate() {
	c.TotalPrice = float64(0)
	recalculatedItems := make([]CartItem, 0)
	for _, item := range c.Items {
		item.Recalculate()
		recalculatedItems = append(recalculatedItems, item)
		c.TotalPrice += item.TotalPrice
	}
	c.Items = recalculatedItems
}


func (c *Cart) IsDeleted() (deleted bool) {
	return c.DeletedAt.Valid && c.DeletedBy.Valid
}

func (c Cart) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ToResponseFormat())
}

func (c *Cart) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(c)
}

func (c *Cart) Update(userID uuid.UUID) (err error){
	c.UpdatedAt = null.TimeFrom(time.Now())
	c.UpdatedBy = nuuid.From(userID)

	c.Recalculate()
	err = c.Validate()
	return
}



func (c Cart) ToResponseFormat() CartResponseFormat {
	resp := CartResponseFormat{
		ID:            	c.ID,
		UserID: 		c.UserID,
		CreatedAt:      c.CreatedAt,
		CreatedBy:     	c.CreatedBy,
		UpdatedAt:      c.UpdatedAt,
		UpdatedBy:     	c.UpdatedBy.Ptr(),
		DeletedAt:      c.DeletedAt,
		DeletedBy:     	c.DeletedBy.Ptr(),
		TotalPrice:     c.TotalPrice,
		Items:         	make([]CartItemResponseFormat, 0),
	}

	for _, item := range c.Items {
		resp.Items = append(resp.Items, item.ToResponseFormat())
	}

	return resp
}

func (c Cart) NewCart(userID uuid.UUID) (newCart Cart, err error){
	cartID, err := uuid.NewV4()
	if err != nil {
		return
	}
	newCart = Cart{
		ID: cartID,
		UserID: userID,
		CreatedAt: time.Now(),
		CreatedBy: userID,
		Items: make([]CartItem, 0),
	}
	return
}

type CartRequestFormat struct {

}

type CartResponseFormat struct {
	ID  			uuid.UUID		`json:"id" validate:"required"`
	UserID  		uuid.UUID		`json:"user_id" validate:"required"`
	CreatedAt     	time.Time   	`json:"created_at" validate:"required"`
	CreatedBy     	uuid.UUID   	`json:"created_by" validate:"required"`
	UpdatedAt     	null.Time   	`json:"updated_at"`
	UpdatedBy     	*uuid.UUID 		`json:"updated_by"`
	DeletedAt     	null.Time  		`json:"deleted_at"`
	DeletedBy     	*uuid.UUID 		`json:"deleted_by"`
	TotalPrice		float64			`json:"total_price"`
	Items           []CartItemResponseFormat `json:"items"`
}


type CartItem struct {
	CartID			uuid.UUID		`db:"cart_id" validate:"required"`
	ProductID		uuid.UUID		`db:"product_id" validate:"required"`
	Quantity		int			`db:"quantity"`
	UnitPrice		float64		`db:"unit_price"`
	TotalPrice		float64		`db:"total_price"`
	ProductStock	int			`db:"stock"`
	CreatedAt		time.Time   `db:"created_at" validate:"required"`
	CreatedBy		uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt		null.Time   `db:"updated_at"`
	UpdatedBy		nuuid.NUUID `db:"updated_by"`
	DeletedAt		null.Time   `db:"deleted_at"`
	DeletedBy		nuuid.NUUID `db:"deleted_by"`
}

func (ci CartItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(ci.ToResponseFormat())
}

func (ci *CartItem) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(ci)
}

func (ci *CartItem) Recalculate() {
	ci.TotalPrice = float64(ci.Quantity) * ci.UnitPrice
}



func (ci CartItem) NewFromRequestFormat(req CartItemRequestFormat, userID uuid.UUID, cartID uuid.UUID) (cartItem CartItem, err error) {
	cartItem = CartItem{
		CartID: cartID,
		ProductID: req.ProductID,
		Quantity: req.Quantity,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	err = ci.Validate()
	return
}

func (ci CartItem) newCartItem(cartID uuid.UUID, req CartItemRequestFormat, userID uuid.UUID) (newCartItem CartItem) {
	newCartItem = CartItem{
		CartID:    cartID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}
	return
}

func (ci *CartItem) Update(cartItem CartItem, userID uuid.UUID) (err error) {
	ci.Quantity += cartItem.Quantity
	if  ci.Quantity < 1 {
		ci.Quantity = 1 
	}
	ci.UnitPrice = cartItem.UnitPrice
	ci.TotalPrice = cartItem.TotalPrice
	ci.UpdatedAt = null.TimeFrom(time.Now())
	ci.UpdatedBy = nuuid.From(userID)

	err = ci.Validate()

	return
}

func (ci *CartItem) ToResponseFormat() CartItemResponseFormat {
	return CartItemResponseFormat{
		CartID: 		ci.CartID,
		ProductID:		ci.ProductID,
		Quantity: 		ci.Quantity,
		UnitPrice:      ci.UnitPrice,
		TotalPrice:  	ci.TotalPrice,	
		CreatedAt:      ci.CreatedAt,
		CreatedBy:     	ci.CreatedBy,
		UpdatedAt:      ci.UpdatedAt,
		UpdatedBy:     	ci.UpdatedBy.Ptr(),
		DeletedAt:      ci.DeletedAt,
		DeletedBy:     	ci.DeletedBy.Ptr(),

	}
}



type CartItemRequestFormat struct {
	ProductID       uuid.UUID    `json:"productID" validate:"required"`
	Quantity		int			`json:"quantity" validate:"required"`
}


type CartItemResponseFormat struct {
	CartID          uuid.UUID `json:"cartID" validate:"required"`
	ProductID       uuid.UUID    `json:"productID" validate:"required"`
	Quantity		int			`json:"quantity"`
	UnitPrice		float64		`json:"unit_price"`
	TotalPrice		float64		`json:"total_price"`
	CreatedAt		time.Time   `json:"created_at" validate:"required"`
	CreatedBy		uuid.UUID   `json:"created_by" validate:"required"`
	UpdatedAt		null.Time   `json:"updated_at"`
	UpdatedBy		*uuid.UUID `json:"updated_by"`
	DeletedAt		null.Time   `json:"deleted_at"`
	DeletedBy		*uuid.UUID `json:"deleted_by"`
}



type CheckoutRequestFormat struct {
	Address 		string 		`json:"address" validate:"required"`
	ProductIDs	[]uuid.UUID    `json:"cart_items" validate:"required"`
}





type CartItemJoin struct {
	CartID			uuid.UUID		`db:"cart_id" validate:"required"`
	ProductID		uuid.UUID		`db:"product_id" validate:"required"`
	Quantity		int			`db:"quantity"`
	UnitPrice		float64		`db:"unit_price"`
	Stock			int			`db:"stock"`
	CreatedAt		time.Time   `db:"created_at" validate:"required"`
	CreatedBy		uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt		null.Time   `db:"updated_at"`
	UpdatedBy		nuuid.NUUID `db:"updated_by"`
	DeletedAt		null.Time   `db:"deleted_at"`
	DeletedBy		nuuid.NUUID `db:"deleted_by"`
}

