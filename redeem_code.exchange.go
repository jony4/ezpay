package ezpay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	APIRedeemCodeExchange = "redeem_code.exchange"
)

// RedeemCodeExchangeRequest redeem_code.exchange
type RedeemCodeExchangeRequest struct {
	CustomerId string `json:"customer_id"` // 艺爪系统用户ID
	Value      string `json:"value"`       // 兑换码
}

// RedeemCodeExchangeResult redeem_code.exchange
type RedeemCodeExchangeResult struct {
	ProjectId string  `json:"project_id"` // 项目ID
	Value     string  `json:"value"`      // 兑换码
	Product   Product `json:"product"`    // 兑换商品信息
}

type Product struct {
	Id   string `json:"id"`   // 商品ID
	Name string `json:"name"` // 商品名称
}

type RedeemCodeExchangeResultError struct {
	ErrorType RedeemCodeExchangeResultErrorType `json:"error"`   // 错误信息
	Message   string                            `json:"message"` // 错误信息
}

func (e *RedeemCodeExchangeResultError) Error() string {
	return fmt.Sprintf("%v:%s", e.ErrorType, e.Message)
}

type RedeemCodeExchangeResultErrorType string

const (
	ErrorUnknown               RedeemCodeExchangeResultErrorType = "Unknown"
	ErrorRedeemCodeAlreadyUsed RedeemCodeExchangeResultErrorType = "RedeemCodeAlreadyUsed"
	ErrorRedeemCodeExhausted   RedeemCodeExchangeResultErrorType = "RedeemCodeExhausted"
	ErrorRedeemCodeExpired     RedeemCodeExchangeResultErrorType = "RedeemCodeExpired"
	ErrorRedeemCodeInvalid     RedeemCodeExchangeResultErrorType = "RedeemCodeInvalid"
)

func (c *Client) RedeemCodeExchange(ctx context.Context, request *RedeemCodeExchangeRequest) (*RedeemCodeExchangeResult, error) {
	resp, err := c.Call(ctx, APIRedeemCodeExchange, request)
	if err != nil && !errors.Is(err, ErrFailedResponse) {
		return nil, err
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	if errors.Is(err, ErrFailedResponse) {
		errResult := &RedeemCodeExchangeResultError{}
		if err := json.Unmarshal(respBody, errResult); err != nil {
			log.Errorf("Unmarshal Error, err=%+v", err)
			return nil, err
		}
		return nil, errResult
	}
	result := &RedeemCodeExchangeResult{}
	if err := json.Unmarshal(respBody, result); err != nil {
		log.Errorf("Unmarshal Error, err=%+v", err)
		return nil, err
	}
	return result, nil
}
