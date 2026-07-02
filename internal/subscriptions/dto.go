package subscriptions

import (
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/sub-tracker-hq/sub-tracker-api/db/generated"
	"github.com/sub-tracker-hq/sub-tracker-api/internal/common"
)

type createRequest struct {
	ProviderKey     *string     `json:"provider_key"`
	Name            string      `json:"name"`
	Category        string      `json:"category"`
	Amount          string      `json:"amount"`
	Currency        string      `json:"currency"`
	BillingCycle    string      `json:"billing_cycle"`
	NextBillingDate pgtype.Date `json:"next_billing_date"`
}

type updateRequest struct {
	Name            string      `json:"name"`
	Category        string      `json:"category"`
	Amount          string      `json:"amount"`
	Currency        string      `json:"currency"`
	BillingCycle    string      `json:"billing_cycle"`
	NextBillingDate pgtype.Date `json:"next_billing_date"`
}

type response struct {
	ID              pgtype.UUID      `json:"id"`
	ProviderKey     *string          `json:"provider_key"`
	Name            string           `json:"name"`
	Category        string           `json:"category"`
	Amount          common.Money     `json:"amount"`
	Currency        string           `json:"currency"`
	BillingCycle    string           `json:"billing_cycle"`
	NextBillingDate pgtype.Date      `json:"next_billing_date"`
	Status          string           `json:"status"`
	Source          string           `json:"source"`
	CreatedAt       common.Timestamp `json:"created_at"`
	UpdatedAt       common.Timestamp `json:"updated_at"`
}

func toResponse(s db.Subscription) response {
	return response{
		ID:              s.ID,
		ProviderKey:     s.ProviderKey,
		Name:            s.Name,
		Category:        s.Category,
		Amount:          common.Money{Numeric: s.Amount},
		Currency:        s.Currency,
		BillingCycle:    s.BillingCycle,
		NextBillingDate: s.NextBillingDate,
		Status:          s.Status,
		Source:          s.Source,
		CreatedAt:       common.Timestamp{Timestamptz: s.CreatedAt},
		UpdatedAt:       common.Timestamp{Timestamptz: s.UpdatedAt},
	}
}
