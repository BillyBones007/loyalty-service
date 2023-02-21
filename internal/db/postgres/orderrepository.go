package postgres

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/BillyBones007/loyalty-service/internal/customerr"
	"github.com/BillyBones007/loyalty-service/internal/db/models"
	"github.com/BillyBones007/loyalty-service/internal/tools/convert"
)

type OrderRepository struct {
	store *Storage
}

// Check order in database. Returned uuid and true (if order is exists) or nil and false (if order is not exists)
func (o *OrderRepository) CheckOrder(order string) (uuid string) {
	var id [16]byte
	q := "SELECT user_id FROM orders WHERE num_order = $1;"
	if err := o.store.Pool.QueryRow(context.TODO(), q, order).Scan(&id); err != nil {
		if errors.Is(err, customerr.ErrNoRows) {
			uuid = ""
			return uuid
		}
		log.Printf("error in function CheckOrder: %s\n", err)
		uuid = ""
		return uuid
	}
	uuid = convert.UUID(id).String()
	return uuid
}

// Insert new order to the database
func (o *OrderRepository) InsertNewOrder(model models.OrderInfo, uuid any) error {
	insertToOrders := `INSERT INTO orders (num_order, time_order, status, user_id)
	VALUES ($1, $2, $3, $4);`
	_, err := o.store.Pool.Exec(context.TODO(), insertToOrders, model.Order, time.Now(), model.Status, uuid)
	if err != nil {
		log.Printf("error insert new order: %s\n", err)
		return err
	}
	return nil
}

// Update order info and balance to the database
func (o *OrderRepository) UpdateOrder(model models.OrderInfo, uuid any) error {
	updateOrder := `UPDATE orders SET status = $1, debet_points = $2
	WHERE num_order = $3;`

	updateBalance := `UPDATE balance SET balance = balance + $1
	WHERE user_id = $2;`
	tx, err := o.store.Pool.Begin(context.TODO())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		}
		tx.Commit(context.TODO())
	}()

	_, err = tx.Exec(context.TODO(), updateOrder, model.Status, model.IntAccrual, model.Order)
	if err != nil {
		log.Printf("error insert order: %s\n", err)
		return err
	}
	_, err = tx.Exec(context.TODO(), updateBalance, model.IntAccrual, uuid)
	if err != nil {
		log.Printf("error update balance: %s\n", err)
		return err
	}
	return nil
}

// Get the current balance for user
func (o *OrderRepository) GetCurrentBalance(model *models.CurrentBalance, uuid any) error {
	q := `SELECT balance, withdraw FROM balance WHERE user_id = $1;`
	var curr int64
	var withdraw int64
	if err := o.store.Pool.QueryRow(context.TODO(), q, uuid).Scan(&curr, &withdraw); err != nil {
		log.Printf("error get current balance: %s\n", err)
		return err
	}
	model.IntCurrent = curr
	model.IntWithdrawn = withdraw
	return nil
}

// Withdrawal points
func (o *OrderRepository) WithdrawalPoints(model models.BalanceWithdraw, uuid any) error {
	currBalanceQuery := `SELECT balance FROM balance WHERE user_id = $1;`
	insertOrderQuery := `INSERT INTO orders (num_order, time_order, status, credit_points, user_id)
	VALUES ($1, $2, $3, $4, $5);`
	updateWithdrawQuery := `UPDATE balance SET withdraw = withdraw + $1, balance = balance - $1 
	WHERE user_id = $2;`
	var currBalance int64
	tx, err := o.store.Pool.Begin(context.TODO())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		}
		tx.Commit(context.TODO())
	}()
	if err = tx.QueryRow(context.TODO(), currBalanceQuery, uuid).Scan(&currBalance); err != nil {
		log.Printf("error withdrawal points: %s\n", err)
		return err
	}
	if currBalance < model.IntSum {
		return customerr.ErrInsufficientFounds
	}
	_, err = tx.Exec(context.TODO(), insertOrderQuery, model.Order, time.Now(), models.Processed, model.IntSum, uuid)
	if err != nil {
		log.Printf("error withdrawal points: %s\n", err)
		return err
	}
	_, err = tx.Exec(context.TODO(), updateWithdrawQuery, model.IntSum, uuid)
	if err != nil {
		log.Printf("error update withdrawal points: %s\n", err)
		return err
	}
	return nil
}

// Get list orders info
func (o *OrderRepository) GetOrdersInfo(uuid any) ([]models.ListOrdersInfo, error) {
	listOrders := make([]models.ListOrdersInfo, 0)
	q := `SELECT num_order, status, debet_points, time_order FROM orders 
	WHERE user_id = $1 ORDER BY time_order;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := o.store.Pool.Query(ctx, q, uuid)
	if err != nil {
		log.Printf("error get order info: %s\n", err)
		return nil, err
	}
	for rows.Next() {
		m := models.ListOrdersInfo{}
		err := rows.Scan(&m.Number, &m.Status, &m.IntAccrual, &m.TimeOrder)
		if err != nil {
			log.Printf("error get order info: %s\n", err)
			return nil, err
		}
		m.Accrual = convert.ConvToFloatBalance(m.IntAccrual)
		listOrders = append(listOrders, m)
	}
	return listOrders, nil
}

// Get list withdrawals info
func (o *OrderRepository) GetWithdrawalsInfo(uuid any) ([]models.Withdrawals, error) {
	listWithdrawals := make([]models.Withdrawals, 0)
	q := `SELECT num_order, credit_points, time_order FROM orders 
	WHERE user_id = $1 ORDER BY time_order DESC;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := o.store.Pool.Query(ctx, q, uuid)
	if err != nil {
		log.Printf("error get order info: %s\n", err)
		return nil, err
	}
	for rows.Next() {
		m := models.Withdrawals{}
		err := rows.Scan(&m.Order, &m.IntSum, &m.TimeOrder)
		if err != nil {
			log.Printf("error get order info: %s\n", err)
			return nil, err
		}
		m.Sum = convert.ConvToFloatBalance(m.IntSum)
		listWithdrawals = append(listWithdrawals, m)
	}
	return listWithdrawals, nil
}

// Select unprocessed orders
func (o *OrderRepository) UnprocessedOrders() ([]models.UnprocessedOrder, error) {
	listOrders := make([]models.UnprocessedOrder, 0)
	var id [16]byte
	q := `SELECT num_order, user_id from orders WHERE status = $1 OR status = $2 OR status = $3;`
	rows, err := o.store.Pool.Query(context.TODO(), q, models.New, models.Processing, models.Registered)
	if err != nil {
		log.Printf("error select unprocessed orders: %s\n", err)
		return nil, err
	}
	for rows.Next() {
		order := models.UnprocessedOrder{}
		err := rows.Scan(&order.Order, &id)
		if err != nil {
			log.Printf("error select unprocessed orders: %s\n", err)
			return nil, err
		}
		order.UUID = convert.UUID(id).String()
		listOrders = append(listOrders, order)
	}
	return listOrders, nil
}
