package repositories

import (
	"database/sql"
	"nabung-emas-api/internal/models"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetPortfolioSummary(userID string) (*models.PortfolioSummary, error) {
	query := `
		SELECT 
			COALESCE(SUM(t.total_price), 0) as total_value,
			COALESCE(SUM(t.weight), 0) as total_weight,
			COUNT(DISTINCT p.id) as total_pockets,
			COUNT(t.id) as total_transactions,
			CASE 
				WHEN SUM(t.weight) > 0 
				THEN SUM(t.total_price) / SUM(t.weight) 
				ELSE 0 
			END as average_price_per_gram
		FROM users u
		LEFT JOIN pockets p ON p.user_id = u.id
		LEFT JOIN transactions t ON t.user_id = u.id
		WHERE u.id = $1
		GROUP BY u.id, t.id
	`

	summary := &models.PortfolioSummary{}
	err := r.db.QueryRow(query, userID).Scan(
		&summary.TotalValue,
		&summary.TotalWeight,
		&summary.TotalPockets,
		&summary.TotalTransactions,
		&summary.AveragePricePerGram,
	)

	if err == sql.ErrNoRows {
		return &models.PortfolioSummary{}, nil
	}

	return summary, err
}

func (r *AnalyticsRepository) GetPocketDistribution(userID string) ([]models.PocketDistribution, error) {
	query := `
		SELECT 
			p.id,
			p.name,
			tp.name,
			tp.color,
			p.aggregate_total_weight,
			p.aggregate_total_price
		FROM pockets p
		LEFT JOIN type_pockets tp ON tp.id = p.type_pocket_id
		WHERE p.user_id = $1 AND p.aggregate_total_weight > 0
		ORDER BY p.aggregate_total_weight DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var distributions []models.PocketDistribution
	var totalWeight float64

	// First pass: collect data and calculate total
	type tempDist struct {
		PocketID        string
		PocketName      string
		TypePocketName  string
		TypePocketColor string
		Weight          float64
		Value           float64
	}
	var tempDistributions []tempDist

	for rows.Next() {
		var td tempDist
		err := rows.Scan(
			&td.PocketID,
			&td.PocketName,
			&td.TypePocketName,
			&td.TypePocketColor,
			&td.Weight,
			&td.Value,
		)
		if err != nil {
			return nil, err
		}
		totalWeight += td.Weight
		tempDistributions = append(tempDistributions, td)
	}

	// Second pass: calculate percentages
	for _, td := range tempDistributions {
		percentage := 0.0
		if totalWeight > 0 {
			percentage = (td.Weight / totalWeight) * 100
		}

		distributions = append(distributions, models.PocketDistribution{
			PocketID:        td.PocketID,
			PocketName:      td.PocketName,
			TypePocketName:  td.TypePocketName,
			TypePocketColor: td.TypePocketColor,
			Weight:          td.Weight,
			Value:           td.Value,
			Percentage:      percentage,
		})
	}

	return distributions, rows.Err()
}

func (r *AnalyticsRepository) GetMonthlyPurchases(userID string, months int, pocketID *string) ([]models.MonthlyPurchaseData, error) {
	query := `
		SELECT 
			TO_CHAR(transaction_date, 'YYYY-MM') as month,
			SUM(weight) as weight,
			SUM(total_price) as amount,
			COUNT(*) as count,
			AVG(price_per_gram) as average_price_per_gram
		FROM transactions
		WHERE user_id = $1
			AND transaction_date >= CURRENT_DATE - INTERVAL '%d months'
	`

	args := []interface{}{userID}
	if pocketID != nil && *pocketID != "" {
		query += " AND pocket_id = $2"
		args = append(args, *pocketID)
	}

	query += `
		GROUP BY TO_CHAR(transaction_date, 'YYYY-MM')
		ORDER BY month DESC
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyData []models.MonthlyPurchaseData
	for rows.Next() {
		var md models.MonthlyPurchaseData
		err := rows.Scan(
			&md.Month,
			&md.Weight,
			&md.Amount,
			&md.Count,
			&md.AveragePricePerGram,
		)
		if err != nil {
			return nil, err
		}
		monthlyData = append(monthlyData, md)
	}

	return monthlyData, rows.Err()
}

func (r *AnalyticsRepository) GetBrandDistribution(userID string) ([]models.BrandDistribution, error) {
	query := `
		SELECT 
			brand,
			SUM(weight) as weight,
			SUM(total_price) as value,
			COUNT(*) as transaction_count
		FROM transactions
		WHERE user_id = $1
		GROUP BY brand
		ORDER BY weight DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var distributions []models.BrandDistribution
	var totalWeight float64

	// First pass: collect data
	type tempBrand struct {
		Brand            string
		Weight           float64
		Value            float64
		TransactionCount int
	}
	var tempBrands []tempBrand

	for rows.Next() {
		var tb tempBrand
		err := rows.Scan(
			&tb.Brand,
			&tb.Weight,
			&tb.Value,
			&tb.TransactionCount,
		)
		if err != nil {
			return nil, err
		}
		totalWeight += tb.Weight
		tempBrands = append(tempBrands, tb)
	}

	// Second pass: calculate percentages
	for _, tb := range tempBrands {
		percentage := 0.0
		if totalWeight > 0 {
			percentage = (tb.Weight / totalWeight) * 100
		}

		distributions = append(distributions, models.BrandDistribution{
			Brand:            tb.Brand,
			Weight:           tb.Weight,
			Value:            tb.Value,
			TransactionCount: tb.TransactionCount,
			Percentage:       percentage,
		})
	}

	return distributions, rows.Err()
}
