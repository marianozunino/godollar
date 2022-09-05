package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/marianozunino/godollar/internal/app/database"
	"github.com/marianozunino/godollar/internal/ine"
	"github.com/marianozunino/godollar/internal/model"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.uber.org/fx"
)

// interface for the ine service
type RateRepository interface {
	GetAll() ([]*model.Rate, error)
	CreateMany(rates []ine.DolarPrice) error
	GetLatest() (*model.Rate, error)
	GetByDate(date string) (*model.Rate, error)
	GetByDateRange(from string, to string) ([]*model.Rate, error)
}

var Module = fx.Options(
	fx.Provide(registerRepo),
)

type rateRepository struct {
	db database.DB
}

// assert that Ine implements IneService
var _ RateRepository = (*rateRepository)(nil)

func registerRepo(db database.DB) RateRepository {
	return &rateRepository{
		db: db,
	}
}

func (r *rateRepository) GetAll() ([]*model.Rate, error) {
	ctx := context.Background()
	rates, err := model.Rates().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *rateRepository) CreateMany(rates []ine.DolarPrice) error {
	fmt.Printf("\nCreating rates: %v", len(rates))
	ctx := context.Background()
	for _, rate := range rates {
		dateString := null.StringFrom(rate.Date.Format("2006-01-02"))

		// create if not exists
		_, err := model.Rates(model.RateWhere.Date.EQ(dateString)).One(ctx, r.db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				var buyPrice null.Float32 = null.Float32From(float32(rate.BuyPrice))
				var sellPrice null.Float32 = null.Float32From(float32(rate.SellPrice))
				var ebrouBuyPrice null.Float32 = null.Float32From(float32(rate.EbrouBuyPrice))
				var ebrouSellPrice null.Float32 = null.Float32From(float32(rate.EbrouSellPrice))

				dbRate := &model.Rate{
					Date: dateString,

					BuyPrice:  buyPrice,
					SellPrice: sellPrice,

					EbrouBuyPrice:  ebrouBuyPrice,
					EbrouSellPrice: ebrouSellPrice,
				}

				if err := dbRate.Insert(ctx, r.db, boil.Infer()); err != nil {
					return err
				}
			} else {
				return err
			}

		}
	}
	return nil
}

func (r *rateRepository) GetLatest() (*model.Rate, error) {
	ctx := context.Background()
	rate, err := model.Rates(qm.OrderBy("date desc")).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func (r *rateRepository) GetByDate(date string) (*model.Rate, error) {
	ctx := context.Background()
	rate, err := model.Rates(model.RateWhere.Date.EQ(
		null.StringFrom(date),
	)).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func (r *rateRepository) GetByDateRange(from string, to string) ([]*model.Rate, error) {
	fmt.Printf("\nGetting rates from %s to %s", from, to)
	ctx := context.Background()

	queryMods := []qm.QueryMod{}

	if from != "" {
		queryMods = append(queryMods, model.RateWhere.Date.GTE(null.StringFrom(from)))
	}

	if to != "" {
		queryMods = append(queryMods, model.RateWhere.Date.LTE(null.StringFrom(to)))
	}

	rates, err := model.Rates(queryMods...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rates, nil
}
