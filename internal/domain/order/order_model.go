package order

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Order struct {
	ID     		uuid.UUID 	`db:"id" validate:"required"`
	UserID 		uuid.UUID  	`db:"user_id" validate:"required"`
	Address 	string		`db:"address" validate:"required"`
	Status 		string		`db:"status" validate:"required"`
	CreatedAt	time.Time   `db:"created_at" validate:"required"`
	CreatedBy	uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt	null.Time   `db:"updated_at"`
	UpdatedBy	nuuid.NUUID `db:"updated_by"`
	DeletedAt	null.Time   `db:"deleted_at"`
	DeletedBy	nuuid.NUUID `db:"deleted_by"`
	TotalPrice	float64		`db:"-"`
	Items		[]OrderItem `db:"-"`
}


func (o *Order) AttachItems(items []OrderItem) Order {
	for _, item := range items {
		if item.OrderID == o.ID {
			o.Items = append(o.Items, item)
		}
	}
	return *o
}

func (o *Order) Recalculate() {
	o.TotalPrice = float64(0)
	recalculatedItems := make([]OrderItem, 0)
	for _, item := range o.Items {
		item.Recalculate()
		recalculatedItems = append(recalculatedItems, item)
		o.TotalPrice += item.TotalPrice
	}
	o.Items = recalculatedItems
}


func (o *Order) IsDeleted() (deleted bool) {
	return o.DeletedAt.Valid && o.DeletedBy.Valid
}

func (o Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ToResponseFormat())
}

func (c *Order) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(c)
}

func (o Order) ToResponseFormat() OrderResponseFormat {
	resp := OrderResponseFormat{
		ID:            	o.ID,
		UserID: 		o.UserID,
		Address:        o.Address,
		Status:         o.Status,
		CreatedAt:      o.CreatedAt,
		CreatedBy:     	o.CreatedBy,
		UpdatedAt:      o.UpdatedAt,
		UpdatedBy:     	o.UpdatedBy.Ptr(),
		DeletedAt:      o.DeletedAt,
		DeletedBy:     	o.DeletedBy.Ptr(),
		TotalPrice: 	o.TotalPrice,
		Items:         	make([]OrderItemResponseFormat, 0),
	}

	for _, item := range o.Items {
		resp.Items = append(resp.Items, item.ToResponseFormat())
	}

	return resp
}

type OrderRequestFormat struct {

}

type OrderResponseFormat struct {
	ID  			uuid.UUID		`json:"id" validate:"required"`
	UserID  		uuid.UUID		`json:"user_id" validate:"required"`
	Address			string			`json:"address"`
	Status 			string 			`json:"status"`
	CreatedAt     	time.Time   	`json:"created_at" validate:"required"`
	CreatedBy     	uuid.UUID   	`json:"created_by" validate:"required"`
	UpdatedAt     	null.Time   	`json:"updated_at"`
	UpdatedBy     	*uuid.UUID 		`json:"updated_by"`
	DeletedAt     	null.Time  		`json:"deleted_at"`
	DeletedBy     	*uuid.UUID 		`json:"deleted_by"`
	TotalPrice		float64			`json:"total_price"`
	Items           []OrderItemResponseFormat `json:"items"`
}

type OrderItem struct {
	OrderID			uuid.UUID		`db:"order_id" validate:"required"`
	ProductID		uuid.UUID		`db:"product_id" validate:"required"`
	Quantity		int			`db:"quantity"`
	UnitPrice		float64		`db:"unit_price"`
	TotalPrice		float64		`db:"-"`
	CreatedAt		time.Time   `db:"created_at" validate:"required"`
	CreatedBy		uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedAt		null.Time   `db:"updated_at"`
	UpdatedBy		nuuid.NUUID `db:"updated_by"`
	DeletedAt		null.Time   `db:"deleted_at"`
	DeletedBy		nuuid.NUUID `db:"deleted_by"`
}


func (o OrderItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ToResponseFormat())
}

func (o *OrderItem) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(o)
}

func (o *OrderItem) Recalculate() {
	o.TotalPrice = float64(o.Quantity) * o.UnitPrice
}


func (oi *OrderItem) ToResponseFormat() OrderItemResponseFormat {
	return OrderItemResponseFormat{
		OrderID: 		oi.OrderID,
		ProductID:		oi.ProductID,
		Quantity: 		oi.Quantity,
		UnitPrice:      oi.UnitPrice,
		TotalPrice:     oi.TotalPrice,	
		CreatedAt:      oi.CreatedAt,
		CreatedBy:     	oi.CreatedBy,
		UpdatedAt:      oi.UpdatedAt,
		UpdatedBy:     	oi.UpdatedBy.Ptr(),
		DeletedAt:      oi.DeletedAt,
		DeletedBy:     	oi.DeletedBy.Ptr(),

	}
}


type OrderItemResponseFormat struct {
	OrderID          uuid.UUID `json:"cartID" validate:"required"`
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














