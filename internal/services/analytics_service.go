package services

import (
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/repositories"
)

type AnalyticsService struct {
	analyticsRepo   *repositories.AnalyticsRepository
	transactionRepo *repositories.TransactionRepository
	pocketRepo      *repositories.PocketRepository
}

func NewAnalyticsService(
	analyticsRepo *repositories.AnalyticsRepository,
	transactionRepo *repositories.TransactionRepository,
	pocketRepo *repositories.PocketRepository,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:   analyticsRepo,
		transactionRepo: transactionRepo,
		pocketRepo:      pocketRepo,
	}
}

func (s *AnalyticsService) GetDashboard(userID string, currentGoldPrice *float64) (*models.DashboardSummary, error) {
	// Get portfolio summary
	portfolio, err := s.analyticsRepo.GetPortfolioSummary(userID)
	if err != nil {
		return nil, err
	}

	// Calculate profit/loss if current price provided
	if currentGoldPrice != nil && *currentGoldPrice > 0 && portfolio.TotalWeight > 0 {
		currentValue := portfolio.TotalWeight * (*currentGoldPrice)
		profitLoss := currentValue - portfolio.TotalValue
		profitLossPercentage := (profitLoss / portfolio.TotalValue) * 100

		portfolio.CurrentGoldPrice = currentGoldPrice
		portfolio.CurrentValue = &currentValue
		portfolio.ProfitLoss = &profitLoss
		portfolio.ProfitLossPercentage = &profitLossPercentage
	}

	// Get recent transactions
	recentTransactions, err := s.transactionRepo.GetRecentTransactions(userID, 5)
	if err != nil {
		return nil, err
	}

	topPockets, err := s.pocketRepo.GetTopPockets(userID)
	if err != nil {
		return nil, err
	}

	return &models.DashboardSummary{
		Portfolio:          *portfolio,
		TopPockets:         topPockets,
		RecentTransactions: recentTransactions,
	}, nil
}

func (s *AnalyticsService) GetPortfolio(userID string, currentGoldPrice *float64) (*models.PortfolioAnalytics, error) {
	// Get portfolio summary
	summary, err := s.analyticsRepo.GetPortfolioSummary(userID)
	if err != nil {
		return nil, err
	}

	// Get distribution
	distribution, err := s.analyticsRepo.GetPocketDistribution(userID)
	if err != nil {
		return nil, err
	}

	analytics := &models.PortfolioAnalytics{
		TotalValue:          summary.TotalValue,
		TotalWeight:         summary.TotalWeight,
		AveragePricePerGram: summary.AveragePricePerGram,
		Distribution:        distribution,
	}

	// Calculate profit/loss if current price provided
	if currentGoldPrice != nil && *currentGoldPrice > 0 && analytics.TotalWeight > 0 {
		currentValue := analytics.TotalWeight * (*currentGoldPrice)
		profitLoss := currentValue - analytics.TotalValue
		profitLossPercentage := (profitLoss / analytics.TotalValue) * 100

		analytics.CurrentMarketPrice = currentGoldPrice
		analytics.CurrentValue = &currentValue
		analytics.ProfitLoss = &profitLoss
		analytics.ProfitLossPercentage = &profitLossPercentage
	}

	return analytics, nil
}

func (s *AnalyticsService) GetMonthlyPurchases(userID string, months int, pocketID *string) (*models.MonthlyPurchaseAnalytics, error) {
	if months < 1 {
		months = 6
	}

	monthlyData, err := s.analyticsRepo.GetMonthlyPurchases(userID, months, pocketID)
	if err != nil {
		return nil, err
	}

	// Calculate averages
	var totalWeight, totalAmount float64
	for _, data := range monthlyData {
		totalWeight += data.Weight
		totalAmount += data.Amount
	}

	avgMonthlyPurchase := 0.0
	if len(monthlyData) > 0 {
		avgMonthlyPurchase = totalAmount / float64(len(monthlyData))
	}

	return &models.MonthlyPurchaseAnalytics{
		MonthlyData:            monthlyData,
		AverageMonthlyPurchase: avgMonthlyPurchase,
		TotalPeriodWeight:      totalWeight,
		TotalPeriodAmount:      totalAmount,
	}, nil
}

func (s *AnalyticsService) GetBrandDistribution(userID string) ([]models.BrandDistribution, error) {
	return s.analyticsRepo.GetBrandDistribution(userID)
}

func (s *AnalyticsService) GetTrends(userID string, period, groupBy string) (*models.TrendAnalytics, error) {
	// TODO: Implement trend analytics with period and groupBy
	// For now, return empty trends
	return &models.TrendAnalytics{
		Trends:  []models.TrendData{},
		Summary: models.TrendsSummary{},
	}, nil
}
