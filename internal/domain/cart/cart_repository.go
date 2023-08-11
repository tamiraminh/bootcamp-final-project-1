package cart

import (
	"database/sql"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var cartQueries = struct {
	selectCart string
	insertCart string
	updateCart string
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
	updateCart: `
	UPDATE atc_cart
	SET
		user_id = :user_id,
		created_at= :created_at,
		created_by= :created_by,
		updated_at= :updated_at,
		updated_by= :updated_by,
		deleted_at= :deleted_at,
		deleted_by= :deleted_by
	WHERE id = :id
	`,
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
	CreateCartByUserID(userID uuid.UUID) (cart Cart, err error)
	UpdateCart(cart Cart) (err error)
	CreateCartItem(cartItem CartItem) (err error)
	UpdateCartItem(cartItem CartItem) (err error)
	DeleteCartItems(cartID uuid.UUID) (err error)
	DeleteCartItem(cartID uuid.UUID, productID uuid.UUID) (err error)
	ResolveCartByUserID(id uuid.UUID) (cart Cart, err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveCartItemByProductID(cartID uuid.UUID, productID uuid.UUID) (cartItem []CartItem, err error)
	ResolveCartItemsByCartID(cartID uuid.UUID) (CartItems []CartItem, err error)
	ResolveCartByID(id uuid.UUID) (cart Cart, err error)
	ResolveCartItemJoinProduct(cartID uuid.UUID, productID uuid.UUID) (cartItem CartItemJoin, err error)
	ResolveCartItemsJoinProduct(cartID uuid.UUID) (cartItem []CartItem, err error)
}

type CartRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCartRepositoryMySQL(db *infras.MySQLConn) *CartRepositoryMySQL {
	s := new(CartRepositoryMySQL)
	s.DB = db
	return s
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

func (r *CartRepositoryMySQL) CreateCartByUserID(userID uuid.UUID) (cart Cart,err error)  {
	_, err = r.ResolveCartByUserID(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	cart, err = cart.NewCart(userID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	err = r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, cart); err != nil {
			e <- err
			return
		}

		e <- nil
	})
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) UpdateCart(cart Cart) (err error) {
	exists, err := r.ExistsByID(cart.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if !exists {
		err = failure.NotFound("Cart")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdateCart(tx, cart); err != nil {
			e <- err
			return
		}

		e <- nil
	})
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

func (r *CartRepositoryMySQL) DeleteCartItems(cartID uuid.UUID) (err error){
	err = r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txDeleteItems(tx, cartID); err != nil {
			e <- err
			return
		}

		e <- nil
	})
	return 
}

func (r *CartRepositoryMySQL) DeleteCartItem(cartID uuid.UUID, productID uuid.UUID) (err error)  {
	err = r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txDeleteItem(tx, cartID, productID); err != nil {
			e <- err
			return
		}

		e <- nil
	})
	return
}

func (r *CartRepositoryMySQL) ResolveCartByUserID(id uuid.UUID) (cart Cart, err error) {
	err = r.DB.Read.Get(
		&cart,
		cartQueries.selectCart+" WHERE user_id = ?",
		id.String())
	if err != nil && err == sql.ErrNoRows {
		log.Info().Msg(err.Error())
		return
	}
	return
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

func (r *CartRepositoryMySQL) ResolveCartItemJoinProduct(cartID uuid.UUID, productID uuid.UUID) (cartItem CartItemJoin, err error)  {
	err = r.DB.Read.Get(
		&cartItem,
		"SELECT cart_id, product_id, quantity, price as unit_price, stock, aci.created_at , aci.created_by, aci.updated_at, aci.updated_by , aci.deleted_at , aci.deleted_by  FROM atc_cart_item aci  JOIN atc_product ap ON aci.product_id = ap.id WHERE cart_id = ? AND product_id = ?",
		cartID.String(), productID.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *CartRepositoryMySQL) ResolveCartItemsJoinProduct(cartID uuid.UUID) (cartItem []CartItem, err error) {
	err = r.DB.Read.Select(
		&cartItem,
		"SELECT cart_id, product_id, quantity, price as unit_price, price*quantity as total_price, stock, aci.created_at , aci.created_by, aci.updated_at, aci.updated_by , aci.deleted_at , aci.deleted_by  FROM atc_cart_item aci  JOIN atc_product ap ON aci.product_id = ap.id WHERE cart_id = ? ",
		cartID.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *CartRepositoryMySQL) ResolveCartByID(id uuid.UUID) (cart Cart, err error)  {
	err = r.DB.Read.Get(
		&cart,
		cartQueries.selectCart+" WHERE id = ?",
		id.String())
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



func (r *CartRepositoryMySQL) txUpdateCart(tx *sqlx.Tx, cart Cart) (err error) {
	stmt, err := tx.PrepareNamed(cartQueries.updateCart)
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


func (r *CartRepositoryMySQL) txDeleteItems(tx *sqlx.Tx, cartID uuid.UUID) (err error) {
	_, err = tx.Exec("DELETE FROM atc_cart_item WHERE cart_id = ?", cartID.String())
	return
}

func (r *CartRepositoryMySQL) txDeleteItem(tx *sqlx.Tx, cartID uuid.UUID, productID uuid.UUID) (err error) {
	_, err = tx.Exec("DELETE FROM atc_cart_item WHERE cart_id = ? AND product_id = ?", cartID.String(), productID.String())
	return
}