package product

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var productQueries = struct {
	selectProduct string
	insertProduct string
}{
	selectProduct: `SELECT * FROM atc_product`,
	insertProduct: `INSERT INTO atc_product (
		id,
		name,
		description,
		category,
		brand,
		stock,
		price,
		created_at,
		created_by,
		updated_at,
		updated_by,
		deleted_at,
		deleted_by
	) VALUES (
		:id,
		:name,
		:description,
		:category,
		:brand,
		:stock,
		:price,
		:created_at,
		:created_by,
		:updated_at,
		:updated_by,
		:deleted_at,
		:deleted_by
	)`,
}

type ProductRepository interface {
	Create(product Product) (err error)
	ExistsByID(id uuid.UUID) (exist bool, err error)
	ResolveProductByID(id uuid.UUID) (product Product, err error)
	ResolveAllProducts(page int, limit int) (products []Product, err error)
}

type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideProductRepositoryMySQL(db *infras.MySQLConn) *ProductRepositoryMySQL {
	s := new(ProductRepositoryMySQL)
	s.DB = db
	return s
}


func (r *ProductRepositoryMySQL) Create(product Product) (err error) {
	exists, err := r.ExistsByID(product.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "product", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, product); err != nil {
			e <- err
			return
		}


		e <- nil
	})
}


func (r *ProductRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(id) FROM atc_product WHERE id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


func (r *ProductRepositoryMySQL) txCreate(tx *sqlx.Tx, product Product) (err error) {
	stmt, err := tx.PrepareNamed(productQueries.insertProduct)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(product)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (r *ProductRepositoryMySQL) ResolveAllProducts(page, limit int) (products []Product, err error) {
	err = r.DB.Read.Select(
		&products,
		productQueries.selectProduct+" ORDER BY name ASC LIMIT ? OFFSET ?", limit, (page)*limit)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


func (r *ProductRepositoryMySQL) ResolveProductByID(id uuid.UUID) (product Product, err error)  {
	err = r.DB.Read.Get(
		&product,
		productQueries.selectProduct+" WHERE id = ?", id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}



	return
}



