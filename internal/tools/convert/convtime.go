package convert

import (
	"time"

	"github.com/BillyBones007/loyalty-service/internal/db/models"
)

// Formats the time into RFC3339 format for models.ListOrderInfo
func TimeInRC3339ForListOrders(listOrders []models.ListOrdersInfo) []models.ListOrdersInfo {
	newList := make([]models.ListOrdersInfo, 0)
	for _, m := range listOrders {
		m.Uploaded = m.TimeOrder.Format(time.RFC3339)
		newList = append(newList, m)
	}
	return newList
}

// Formats the time into RFC3339 format for models.Withdrawals
func TimeInRC3339ForWithdrawals(listWithdrawals []models.Withdrawals) []models.Withdrawals {
	newList := make([]models.Withdrawals, 0)
	for _, m := range listWithdrawals {
		m.Processed = m.TimeOrder.Format(time.RFC3339)
		newList = append(newList, m)
	}
	return newList
}
