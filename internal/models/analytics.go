package models

type PortfolioAnalytics struct {
	TotalValue           float64                  `json:"total_value"`
	TotalWeight          float64                  `json:"total_weight"`
	AveragePricePerGram  float64                  `json:"average_price_per_gram"`
	CurrentMarketPrice   *float64                 `json:"current_market_price,omitempty"`
	CurrentValue         *float64                 `json:"current_value,omitempty"`
	ProfitLoss           *float64                 `json:"profit_loss,omitempty"`
	ProfitLossPercentage *float64                 `json:"profit_loss_percentage,omitempty"`
	Distribution         []PocketDistribution     `json:"distribution"`
}

type PocketDistribution struct {
	PocketID        string  `json:"pocket_id"`
	PocketName      string  `json:"pocket_name"`
	TypePocketName  string  `json:"type_pocket_name"`
	TypePocketColor string  `json:"type_pocket_color"`
	Weight          float64 `json:"weight"`
	Value           float64 `json:"value"`
	Percentage      float64 `json:"percentage"`
}

type MonthlyPurchaseData struct {
	Month              string  `json:"month"`
	Weight             float64 `json:"weight"`
	Amount             float64 `json:"amount"`
	Count              int     `json:"count"`
	AveragePricePerGram float64 `json:"average_price_per_gram"`
}

type MonthlyPurchaseAnalytics struct {
	MonthlyData            []MonthlyPurchaseData `json:"monthly_data"`
	AverageMonthlyPurchase float64               `json:"average_monthly_purchase"`
	TotalPeriodWeight      float64               `json:"total_period_weight"`
	TotalPeriodAmount      float64               `json:"total_period_amount"`
}

type BrandDistribution struct {
	Brand            string  `json:"brand"`
	Weight           float64 `json:"weight"`
	Value            float64 `json:"value"`
	TransactionCount int     `json:"transaction_count"`
	Percentage       float64 `json:"percentage"`
}

type TrendData struct {
	Period              string  `json:"period"`
	TotalWeight         float64 `json:"total_weight"`
	TotalValue          float64 `json:"total_value"`
	TransactionCount    int     `json:"transaction_count"`
	AveragePricePerGram float64 `json:"average_price_per_gram"`
}

type TrendAnalytics struct {
	Trends  []TrendData   `json:"trends"`
	Summary TrendsSummary `json:"summary"`
}

type TrendsSummary struct {
	TotalWeight          float64 `json:"total_weight"`
	TotalValue           float64 `json:"total_value"`
	TransactionCount     int     `json:"transaction_count"`
	AveragePricePerGram  float64 `json:"average_price_per_gram"`
	LowestPricePerGram   float64 `json:"lowest_price_per_gram"`
	HighestPricePerGram  float64 `json:"highest_price_per_gram"`
}

type DashboardSummary struct {
	Portfolio           PortfolioSummary `json:"portfolio"`
	RecentTransactions  []Transaction    `json:"recent_transactions"`
}

type PortfolioSummary struct {
	TotalValue           float64  `json:"total_value"`
	TotalWeight          float64  `json:"total_weight"`
	TotalPockets         int      `json:"total_pockets"`
	TotalTransactions    int      `json:"total_transactions"`
	AveragePricePerGram  float64  `json:"average_price_per_gram"`
	CurrentGoldPrice     *float64 `json:"current_gold_price,omitempty"`
	CurrentValue         *float64 `json:"current_value,omitempty"`
	ProfitLoss           *float64 `json:"profit_loss,omitempty"`
	ProfitLossPercentage *float64 `json:"profit_loss_percentage,omitempty"`
}
