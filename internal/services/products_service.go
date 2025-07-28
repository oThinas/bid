package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oThinas/bid/internal/store/pg"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pg.New(pool),
	}
}

func (ps *ProductService) CreateProduct(
	ctx context.Context,
	sellerID uuid.UUID,
	name, description string,
	basePrice float64,
	auctionEnd time.Time,
) (uuid.UUID, error) {
	id, err := ps.queries.CreateProduct(ctx, pg.CreateProductParams{
		SellerID:    sellerID,
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		AuctionEnd:  auctionEnd,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (ps *ProductService) GetProductByID(ctx context.Context, productID uuid.UUID) (pg.Product, error) {
	product, err := ps.queries.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pg.Product{}, ErrProductNotFound
		}

		return pg.Product{}, err
	}

	return product, nil
}
