package metrics

import (
	"capitalExporter/account" // Replace with the actual path to your account package
	"github.com/prometheus/client_golang/prometheus"
)

var (
	balanceMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "account_balance",
			Help: "Account balance",
		},
		[]string{"accountID"},
	)
	depositMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "account_deposit",
			Help: "Account deposit",
		},
		[]string{"accountID"},
	)
	profitLossMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "account_profit_loss",
			Help: "Account profit/loss",
		},
		[]string{"accountID"},
	)
	availableMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "account_available",
			Help: "Account available",
		},
		[]string{"accountID"},
	)
)

func init() {
	prometheus.MustRegister(balanceMetric, depositMetric, profitLossMetric, availableMetric)
}

func UpdateMetrics(sessionTokens account.SessionResponse) error {
	// Call refreshTokenIfNeeded from the accountInfo package
	sessionTokens, err := account.RefreshTokenIfNeeded(sessionTokens)
	if err != nil {
		return err
	}

	// Call getDetails from the account package
	accounts, err := account.GetDetails(sessionTokens.CST, sessionTokens.XSecurityToken)
	if err != nil {
		return err
	}

	// Iterate through accounts and update metrics
	for _, accountInfo := range accounts {
		balanceMetric.WithLabelValues(accountInfo.AccountID).Set(accountInfo.Balance.Balance)
		depositMetric.WithLabelValues(accountInfo.AccountID).Set(accountInfo.Balance.Deposit)
		profitLossMetric.WithLabelValues(accountInfo.AccountID).Set(accountInfo.Balance.ProfitLoss)
		availableMetric.WithLabelValues(accountInfo.AccountID).Set(accountInfo.Balance.Available)
	}

	return nil
}
