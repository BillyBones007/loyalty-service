package db

import "github.com/BillyBones007/loyalty-service/internal/db/models"

type UniversalUserRepository interface {
	Create(model *models.User) error
	UserIsExists(model *models.User) bool
	FindByLogin(model *models.User) error
}

type UniversalOrderRepository interface {
	CheckOrder(order string) (uuid string)
	InsertNewOrder(model models.OrderInfo, uuid any) error
	UpdateOrder(model models.OrderInfo, uuid any) error
	GetCurrentBalance(model *models.CurrentBalance, uuid any) error
	GetOrdersInfo(uuid any) ([]models.ListOrdersInfo, error)
	GetWithdrawalsInfo(uuid any) ([]models.Withdrawals, error)
	WithdrawalPoints(model models.BalanceWithdraw, uuid any) error
	UnprocessedOrders() ([]models.UnprocessedOrder, error)
}
