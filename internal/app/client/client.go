package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/BillyBones007/loyalty-service/internal/customerr"
	"github.com/BillyBones007/loyalty-service/internal/db"
	"github.com/BillyBones007/loyalty-service/internal/db/models"
	"github.com/BillyBones007/loyalty-service/internal/tools/convert"
)

// Client
type AccrualClient struct {
	Client      *http.Client
	Storage     db.Store
	AddrAccrual string
}

// New client
func NewAccrualClient(storage db.Store, addrAccrual string) *AccrualClient {
	return &AccrualClient{Client: &http.Client{}, Storage: storage, AddrAccrual: addrAccrual}
}

// Run client
func (c *AccrualClient) Run() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	errCh := make(chan struct{})
	stopCh := make(chan struct{})
	orders := make(chan models.UnprocessedOrder, 50)
	results := make(chan models.OrderInfo)

	for i := 0; i < cap(orders); i++ {
		go c.WorkerOrderInfo(orders, results, errCh, stopCh)
	}

	for {
		select {
		case <-errCh:
			stopCh <- struct{}{}
		case <-ticker.C:
			c.CheckUnprocessedOrders(orders)
		}

	}
}

// Check unprocessed orders
func (c *AccrualClient) CheckUnprocessedOrders(orders chan models.UnprocessedOrder) {
	listOrders, err := c.Storage.Order().UnprocessedOrders()
	if err != nil {
		if errors.Is(err, customerr.ErrNoRows) {
			return
		}
		fmt.Printf("error from client: %s\n", err)
		return
	}
	for _, order := range listOrders {
		orders <- order
	}
}

// GET info to the accrual system
func (c *AccrualClient) WorkerOrderInfo(orderCh chan models.UnprocessedOrder, resCh chan models.OrderInfo, errCh chan struct{}, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			time.Sleep(3 * time.Second)
		case order := <-orderCh:
			endpoint := fmt.Sprintf("%s/api/orders/%s", c.AddrAccrual, order.Order)
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)

			if err != nil {
				log.Printf("error from worker: %s\n", err)
				continue
			}

			resp, err := c.Client.Do(req)

			if err != nil {
				log.Printf("error from worker: %s\n", err)
				continue
			}

			sCode := resp.StatusCode

			if sCode == http.StatusTooManyRequests {
				errCh <- struct{}{}

			} else if sCode == http.StatusNoContent {
				orderInfo := models.OrderInfo{Order: order.Order, Status: models.Invalid, UUID: order.UUID}
				c.Storage.Order().UpdateOrder(orderInfo, orderInfo.UUID)
				continue

			} else if sCode == http.StatusInternalServerError {
				continue

			} else if sCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("error from worker: %s\n", err)
					resp.Body.Close()
					continue
				}
				resp.Body.Close()
				orderInfo := models.OrderInfo{UUID: order.UUID}
				err = json.Unmarshal(body, &orderInfo)
				if err != nil {
					log.Printf("error from worker: %s\n", err)
					continue
				}
				if orderInfo.Accrual != 0 {
					orderInfo.IntAccrual = convert.ConvToIntBalance(orderInfo.Accrual)
					c.Storage.Order().UpdateOrder(orderInfo, orderInfo.UUID)
				}
			}
		}
	}
}
