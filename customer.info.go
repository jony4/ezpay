package ezpay

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	APICustomerInfo = "customer.info"
)

type EZCustomerInfoRequest struct {
	PaywallID         string   `json:"paywall_id"`                    // 付费界面ID，ID和别名传一个即可
	PaywallAlias      string   `json:"paywall_alias,omitempty"`       // 付费界面别名，ID和别名传一个即可
	ExtraPaywallID    string   `json:"extra_paywall_id,omitempty"`    // 可选，额外的付费界面ID
	ExtraPaywallAlias string   `json:"extra_paywall_alias,omitempty"` // 可选，额外的付费界面别名
	ClientOpenId      string   `json:"client_openid,omitempty"`       // 可选，微信支付AppID:OpenID，中间冒号分隔
	Customer          Customer `json:"customer"`
	IncludeBalance    bool     `json:"include_balance"` // 是否返回用户余额
}

type Customer struct {
	ExternalID        string `json:"external_id"`                   // 商户系统用户ID
	Nickname          string `json:"nickname,omitempty"`            // 商户系统用户用户名/昵称，可选
	ExternalDtCreated string `json:"external_dt_created,omitempty"` // 商户系统用户创建时间，可选
	Promoter          string `json:"promoter,omitempty"`            // 推广员参数，可选
}

type EZCustomerInfo struct {
	ID                string        `json:"id"`                  // 艺爪系统用户ID
	ExternalID        string        `json:"external_id"`         // 商户系统用户ID
	ExternalDtCreated time.Time     `json:"external_dt_created"` // 商户系统用户创建时间
	Nickname          string        `json:"nickname"`
	BalanceS          []BalanceInfo `json:"balance_s"`  // 余额信息
	HomeLink          HomeLink      `json:"home_link"`  // 对应 paywall_id 和 paywall_alias
	ExtraLink         ExtraLink     `json:"extra_link"` // 对应 extra_paywall_id 和 extra_paywall_alias
	DtCreated         time.Time     `json:"dt_created"`
	DtUpdated         time.Time     `json:"dt_updated"`
}

type BalanceInfo struct {
	Equity            Equity `json:"equity"`              // 权益信息
	Balance           int64  `json:"balance"`             // 到期时间或积分数量
	BalanceText       string `json:"balance_text"`        // 余额/到期时间字符串
	IsBalanceInfinite bool   `json:"is_balance_infinite"` // 是否永久或无限
	IsBalanceUsable   bool   `json:"is_balance_usable"`   // 是否有余额/会员是否有效
	HasCharged        bool   `json:"has_charged"`         // 是否有效充值或兑换过
	HasInitial        bool   `json:"has_initial"`         // 是否用过试用会员
	BalanceLogID      string `json:"balance_log_id"`      // 最后一次余额更新日志ID
}

type Equity struct {
	ID    string `json:"id"`    // 权益ID
	Name  string `json:"name"`  // 权益名称
	Alias string `json:"alias"` // 权益别名
}

type HomeLink struct {
	URL         string `json:"url"`          // 付费页面链接，含用户登录凭证
	PromoterURL string `json:"promoter_url"` // 推广员链接，含用户登录凭证
}

// ExtraLink 暂时用不到
type ExtraLink struct {
	URL         string `json:"url"`          // 付费页面链接，含用户登录凭证
	PromoterURL string `json:"promoter_url"` // 推广员链接，含用户登录凭证
}

func (e *EZCustomerInfo) String() string {
	es, _ := json.Marshal(e)
	return string(es)
}

func (c *Client) CustomerInfo(ctx context.Context, request *EZCustomerInfoRequest) (*EZCustomerInfo, error) {
	resp, err := c.call(ctx, APICustomerInfo, request)
	if err != nil {
		return nil, err
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	log.Debugf("respBody=%s", string(respBody))

	ezUserInfo := &EZCustomerInfo{}
	if err := json.Unmarshal(respBody, ezUserInfo); err != nil {
		log.Errorf("Unmarshal Error, err=%+v", err)
		return nil, err
	}
	return ezUserInfo, nil
}
