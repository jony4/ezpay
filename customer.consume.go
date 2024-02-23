package ezpay

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

const (
	APICustomerConsume = "customer.consume"
)

// ConsumeRequest customer.consume
type ConsumeRequest struct {
	CustomerID   string     `json:"customer_id"`    // 艺爪系统用户ID
	EquityID     string     `json:"equity_id"`      // 权益ID
	Amount       int64      `json:"amount"`         // 数量
	Title        string     `json:"title"`          // 标题
	ChangeType   ChangeType `json:"change_type"`    // 变更类型枚举值
	BalanceLogID string     `json:"balance_log_id"` // 最后一次余额更新日志ID
}

type ChangeType int

const (
	ChangeTypeCONSUMABLE_EXPENSE ChangeType = 16
	ChangeTypeCONSUMABLE_INCOME  ChangeType = 17
)

// ConsumeResult customer.consume
type ConsumeResult struct {
	EquityID          string `json:"equity_id"`           // 权益ID
	Balance           int64  `json:"balance"`             // 到期时间或积分数量
	IsBalanceInfinite bool   `json:"is_balance_infinite"` // 是否永久或无限
	BalanceLogID      string `json:"balance_log_id"`      // 最后一次余额更新日志ID
}

func (c *Client) Consume(ctx context.Context, request *ConsumeRequest) (*ConsumeResult, error) {
	resp, err := c.Call(ctx, APICustomerConsume, request)
	if err != nil {
		return nil, err
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	eq := &ConsumeResult{}
	if err := json.Unmarshal(respBody, eq); err != nil {
		log.Errorf("Unmarshal Error, err=%+v", err)
		return nil, err
	}
	return eq, nil
}
