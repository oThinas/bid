package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oThinas/bid/internal/store/pg"
)

type BidsService struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
}

func NewBidsService(pool *pgxpool.Pool) BidsService {
	return BidsService{
		pool:    pool,
		queries: pg.New(pool),
	}
}

func (bs *BidsService) PlaceBid(ctx context.Context, productID, bidderID uuid.UUID, amount float64) (pg.Bid, error) {
	product, err := bs.queries.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pg.Bid{}, err
		}
	}

	highestBid, err := bs.queries.GetHighestBidByProductID(ctx, productID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pg.Bid{}, err
		}
	}

	if product.BasePrice >= amount || highestBid.Amount >= amount {
		return pg.Bid{}, ErrBidAmountTooLow
	}

	highestBid, err = bs.queries.CreateBid(ctx, pg.CreateBidParams{
		ProductID: productID,
		BidderID:  bidderID,
		Amount:    amount,
	})
	if err != nil {
		return pg.Bid{}, err
	}

	return highestBid, nil
}
