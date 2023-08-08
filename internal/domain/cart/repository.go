package cart

import (
	"database/sql"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var cartQueries = struct {
	selectCart string
	insertCart string
	selectCartItem string
	insertCartItem string
	updateCartItem string
}{
	selectCart: `SELECT * FROM atc_cart`,
	insertCart: `INSERT INTO atc_cart (
		id,
		user_id,
		created_at,
		created_by,
		updated_at,
		updated_by,
		deleted_at,
		deleted_by
	) VALUES (
		:id,
		:user_id,
		:created_at,
		:created_by,
		:updated_at,
		:updated_by,
		:deleted_at,
		:deleted_by
	)`,
	selectCartItem: `SELECT * FROM atc_cart_item`,
	insertCartItem: ` INSERT INTO atc_cart_item (
		cart_id,
		product_id,
		quantity,
		created_at,
		created_by,
		updated_at,
		updated_by,
		deleted_at,
		deleted_by
	) VALUES (
		:cart_id,
		:product_id,
		:quantity,
		:created_at,
		:created_by,
		:updated_at,
		:updated_by,
		:deleted_at,
		:deleted_by
	)
	`,
	updateCartItem: `UPDATE atc_cart_item
	SET
		quantity= :quantity,
		created_at= :created_at,
		created_by= :created_by,
		updated_at= :updated_at,
		updated_by= :updated_by,
		deleted_at= :deleted_at,
		deleted_by= :deleted_by
	WHERE cart_id =:cart_id and product_id =:product_id
	` ,
}

type CartRepository interface {
	CreateCart(cart Cart) (err error)
	CreateCartItem(cartItem CartItem) (err error)
	UpdateCartItem(cartItem CartItem) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveCartByUserID(id uuid.UUID) (cart Cart, err error)
	ResolveCartItemByProductID(cartID uuid.UUID, productID uuid.UUID) (cartItem []CartItem, err error)
	ResolveCartItemsByCartID(cartID uuid.UUID) (CartItems []CartItem, err error)
}

type CartRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCartRepositoryMySQL(db *infras.MySQLConn) *CartRepositoryMySQL {
	s := new(CartRepositoryMySQL)
	s.DB = db
	return s
}

func (r *CartRepositoryMySQL) ResolveCartByUserID(id uuid.UUID) (cart Cart, err error) {
	err = r.DB.Read.Get(
		&cart,
		cartQueries.selectCart+" WHERE user_id = ?",
		id.String())
	if err != nil && err == sql.ErrNoRows {
		logger.ErrorWithStack(err)
		return
	}
	return
}


func (r *CartRepositoryMySQL) CreateCart(cart Cart) (err error) {
	exists, err := r.ExistsByID(cart.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "foo", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, cart); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

func (r *CartRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM atc_cart WHERE id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


func (r *CartRepositoryMySQL) CreateCartItem(cartItem CartItem) (err error) {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreateItem(tx, cartItem); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

func (r *CartRepositoryMySQL) UpdateCartItem(cartItem CartItem) (err error) {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdateItem(tx, cartItem); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}



func (r *CartRepositoryMySQL) ResolveCartItemByProductID(cartID uuid.UUID, productID uuid.UUID) (cartItem []CartItem, err error) {
	err = r.DB.Read.Select(
		&cartItem,
		cartQueries.selectCartItem+" WHERE product_id = ? and cart_id = ? ",
		productID.String(), cartID.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (r *CartRepositoryMySQL) ResolveCartItemsByCartID(cartID uuid.UUID) (cartItems []CartItem, err error) {
	err = r.DB.Read.Select(
		&cartItems,
		cartQueries.selectCartItem+" WHERE cart_id = ?",
		cartID.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}




func (r *CartRepositoryMySQL) txCreate(tx *sqlx.Tx, cart Cart) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.insertCart)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cart)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


func (r *CartRepositoryMySQL) txCreateItem(tx *sqlx.Tx, cartItem CartItem) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.insertCartItem)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cartItem)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


func (r *CartRepositoryMySQL) txUpdateItem(tx *sqlx.Tx, cartItem CartItem) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.updateCartItem)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cartItem)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
