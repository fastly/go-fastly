package fastly

import (
	"fmt"
	"strconv"
	"time"
)

// Billing is the top-level representation of a billing response from the Fastly
// API.
type Billing struct {
	EndTime   *time.Time     `mapstructure:"end_time"`
	InvoiceID *string        `mapstructure:"invoice_id"`
	StartTime *time.Time     `mapstructure:"start_time"`
	Status    *BillingStatus `mapstructure:"status"`
	Total     *BillingTotal  `mapstructure:"total"`
}

// BillingStatus is a representation of the status of the bill from the Fastly
// API.
type BillingStatus struct {
	InvoiceID *string    `mapstructure:"invoice_id"`
	SentAt    *time.Time `mapstructure:"sent_at"`
	Status    *string    `mapstructure:"status"`
}

// BillingTotal is a representation of the status of the usage for this bill
// from the Fastly API.
type BillingTotal struct {
	Bandwidth          *float64        `mapstructure:"bandwidth"`
	BandwidthCost      *float64        `mapstructure:"bandwidth_cost"`
	Cost               *float64        `mapstructure:"cost"`
	CostBeforeDiscount *float64        `mapstructure:"cost_before_discount"`
	Discount           *float64        `mapstructure:"discount"`
	Extras             []*BillingExtra `mapstructure:"extras"`
	ExtrasCost         *float64        `mapstructure:"extras_cost"`
	IncurredCost       *float64        `mapstructure:"incurred_cost"`
	Overage            *float64        `mapstructure:"overage"`
	PlanCode           *string         `mapstructure:"plan_code"`
	PlanMinimum        *string         `mapstructure:"plan_minimum"`
	PlanName           *string         `mapstructure:"plan_name"`
	Requests           *uint64         `mapstructure:"requests"`
	RequestsCost       *float64        `mapstructure:"requests_cost"`
	Terms              *string         `mapstructure:"terms"`
}

// BillingExtra is a representation of extras (such as SSL addons) from the
// Fastly API.
type BillingExtra struct {
	Name      *string  `mapstructure:"name"`
	Recurring *float64 `mapstructure:"recurring"`
	Setup     *float64 `mapstructure:"setup"`
}

// GetBillingInput is used as input to the GetBilling function.
type GetBillingInput struct {
	// Month is a 2-digit month (required).
	Month uint8
	// Year is a 4-digit year (required).
	Year uint16
}

// GetBilling returns the billing information for the current account.
func (c *Client) GetBilling(i *GetBillingInput) (*Billing, error) {
	if i.Year == 0 {
		return nil, ErrMissingYear
	}
	if i.Month == 0 {
		return nil, ErrMissingMonth
	}

	path := ToSafeURL("billing", "year", strconv.Itoa(int(i.Year)), "month", fmt.Sprintf("%02d", i.Month))

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Billing
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}
