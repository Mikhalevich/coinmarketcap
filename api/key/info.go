package key

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mikhalevich/coinmarketcap"
	"github.com/Mikhalevich/coinmarketcap/api/types"
)

const (
	keyEndpoint = "/v1/key/info"
)

type KeyResponse struct {
	Data   KeyData      `json:"data"`
	Status types.Status `json:"status"`
}

type KeyData struct {
	Plan  KeyPlan  `json:"plan"`
	Usage KeyUsage `json:"usage"`
}

type KeyPlan struct {
	CreditLimitMonthly               float64   `json:"credit_limit_monthly"`
	CreditLimitMonthlyReset          string    `json:"credit_limit_monthly_reset"`
	CreditLimitMonthlyResetTimestamp time.Time `json:"credit_limit_monthly_reset_timestamp"`
	RateLimitMinute                  float64   `json:"rate_limit_minute"`
}

type KeyUsage struct {
	CurrentMinute KeyUsageMinute  `json:"current_minute"`
	CurrentDay    KeyUsageCredits `json:"current_day"`
	CurrentMonth  KeyUsageCredits `json:"current_month"`
}

type KeyUsageMinute struct {
	RequestsMade float64 `json:"requests_made"`
	RequestsLeft float64 `json:"requests_left"`
}

type KeyUsageCredits struct {
	CreditsUsed float64 `json:"credits_used"`
	CreditsLeft float64 `json:"credits_left"`
}

// Info returns API key details and usage stats.
// This endpoint can be used to programmatically monitor your key usage compared to the rate limit
// and daily/monthly credit limits available to your API plan.
// You may use the Developer Portal's account dashboard as an alternative to this endpoint.
// https://coinmarketcap.com/api/documentation/v1/#operation/getV1KeyInfo
func (k *Key) Info(ctx context.Context) (*KeyResponse, error) {
	var response KeyResponse

	if err := k.executor.Get(
		ctx,
		keyEndpoint,
		func(req *http.Request) error {
			return nil
		},
		&response,
	); err != nil {
		return nil, fmt.Errorf("execute get request: %w", err)
	}

	if response.Status.IsError() {
		return nil, coinmarketcap.NewError(response.Status.ErrorCode, response.Status.ErrorMessage)
	}

	return &response, nil
}
