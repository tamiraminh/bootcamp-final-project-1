package order

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var orderQueries = struct {
	insertOrder 	string
	insertOrderItem string
} {
	insertOrder: `INSERT INTO atc_order (
		id,
		user_id,
		address,
		status,
		created_at,
		created_by,
		updated_at,
		updated_by,
		deleted_at,
		deleted_by
	) VALUES (
		:id,
		:user_id,
		:address,
		:status,
		:created_at,
		:created_by,
		:updated_at,
		:updated_by,
		:deleted_at,
		:deleted_by
	)`,
	insertOrderItem: `
	INSERT INTO atc_order_item (
		order_id,
		product_id,
		quantity,
		unit_price,
		created_at,
		created_by,
		updated_at,
		updated_by,
		deleted_at,
		deleted_by
	) VALUES (
		:order_id,
		:product_id,
		:quantity,
		:unit_price,
		:created_at,
		:created_by,
		:updated_at,
		:updated_by,
		:deleted_at,
		:deleted_by
	)
	`,
}

type OrderRepository interface {
	CreateOrder(order Order) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	CreateOrderItem(oi OrderItem) (err error)

}

type OrderRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideOrderRepositoryMySQL(db *infras.MySQLConn) *OrderRepositoryMySQL  {
	s := new(OrderRepositoryMySQL)
	s.DB = db
	return s
}

func (r *OrderRepositoryMySQL) CreateOrder(order Order) (err error) {
	exists, err := r.ExistsByID(order.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "order", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, order); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

func (r *OrderRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM atc_order WHERE id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}
	
	return
}

func (r *OrderRepositoryMySQL) CreateOrderItem(oi OrderItem) (err error)   {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreateItem(tx, oi); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

func (r *OrderRepositoryMySQL) txCreate(tx *sqlx.Tx, order Order) (err error) {
	stmt, err := tx.PrepareNamed(orderQueries.insertOrder)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(order)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (r *OrderRepositoryMySQL) txCreateItem(tx *sqlx.Tx, oi OrderItem) (err error) {
	stmt, err := tx.PrepareNamed(orderQueries.insertOrderItem)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(oi)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


