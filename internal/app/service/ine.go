package service

import (
	"github.com/marianozunino/godollar/internal/app/dto"
	"github.com/marianozunino/godollar/internal/app/repository"
	"github.com/marianozunino/godollar/internal/ine"
	"go.uber.org/fx"
)

// interface for the ine service
type IneService interface {
	Populate() error
	GetAll() ([]dto.RateDto, error)
	GetToday() (dto.RateDto, error)
	GetByDate(date string) (dto.RateDto, error)
	GetByDateRange(from string, to string) ([]dto.RateDto, error)
}

var Module = fx.Options(
	fx.Provide(registerService),
)

func registerService(rateRepo repository.RateRepository) IneService {
	return &ineService{
		rateRepo: rateRepo,
	}
}

type ineService struct {
	rateRepo repository.RateRepository
}

// assert that Ine implements IneService
var _ IneService = (*ineService)(nil)

func (i *ineService) Populate() error {
	prices, error := ine.GetIneDollarPrices()
	if error != nil {
		return error
	}
	i.rateRepo.CreateMany(prices)
	return nil
}

func (i *ineService) GetAll() ([]dto.RateDto, error) {
	rates, err := i.rateRepo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]dto.RateDto, len(rates))

	for i, rate := range rates {
		result[i] = dto.RateDto{
			Date:           rate.Date.String,
			BuyPrice:       float64(rate.BuyPrice.Float32),
			SellPrice:      float64(rate.SellPrice.Float32),
			EbourSellPrice: float64(rate.EbrouSellPrice.Float32),
			EbourBuyPrice:  float64(rate.EbrouBuyPrice.Float32),
		}
	}
	return result, nil
}

func (i *ineService) GetToday() (dto.RateDto, error) {
	rate, err := i.rateRepo.GetLatest()
	if err != nil {
		return dto.RateDto{}, err
	}
	return dto.RateDto{
		Date:           rate.Date.String,
		BuyPrice:       float64(rate.BuyPrice.Float32),
		SellPrice:      float64(rate.SellPrice.Float32),
		EbourSellPrice: float64(rate.EbrouSellPrice.Float32),
		EbourBuyPrice:  float64(rate.EbrouBuyPrice.Float32),
	}, nil
}

func (i *ineService) GetByDate(date string) (dto.RateDto, error) {
	rate, err := i.rateRepo.GetByDate(date)
	if err != nil {
		return dto.RateDto{}, err
	}
	return dto.RateDto{
		Date:           rate.Date.String,
		BuyPrice:       float64(rate.BuyPrice.Float32),
		SellPrice:      float64(rate.SellPrice.Float32),
		EbourSellPrice: float64(rate.EbrouSellPrice.Float32),
		EbourBuyPrice:  float64(rate.EbrouBuyPrice.Float32),
	}, nil
}

func (i *ineService) GetByDateRange(from string, to string) ([]dto.RateDto, error) {
	rates, err := i.rateRepo.GetByDateRange(from, to)
	if err != nil {
		return nil, err
	}
	result := make([]dto.RateDto, len(rates))

	for i, rate := range rates {
		result[i] = dto.RateDto{
			Date:           rate.Date.String,
			BuyPrice:       float64(rate.BuyPrice.Float32),
			SellPrice:      float64(rate.SellPrice.Float32),
			EbourSellPrice: float64(rate.EbrouSellPrice.Float32),
			EbourBuyPrice:  float64(rate.EbrouBuyPrice.Float32),
		}
	}
	return result, nil
}
