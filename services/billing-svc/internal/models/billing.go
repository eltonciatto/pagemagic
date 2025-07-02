package models

import (
	"time"
)

type User struct {
	ID               string    `json:"id" db:"id"`
	Email            string    `json:"email" db:"email"`
	StripeCustomerID string    `json:"stripe_customer_id" db:"stripe_customer_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type Subscription struct {
	ID                   string    `json:"id" db:"id"`
	UserID               string    `json:"user_id" db:"user_id"`
	StripeSubscriptionID string    `json:"stripe_subscription_id" db:"stripe_subscription_id"`
	Status               string    `json:"status" db:"status"`
	PlanID               string    `json:"plan_id" db:"plan_id"`
	CurrentPeriodStart   time.Time `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd     time.Time `json:"current_period_end" db:"current_period_end"`
	CancelAtPeriodEnd    bool      `json:"cancel_at_period_end" db:"cancel_at_period_end"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type Invoice struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	StripeInvoiceID string     `json:"stripe_invoice_id" db:"stripe_invoice_id"`
	SubscriptionID  string     `json:"subscription_id" db:"subscription_id"`
	Status          string     `json:"status" db:"status"`
	Amount          int64      `json:"amount" db:"amount"`
	Currency        string     `json:"currency" db:"currency"`
	DueDate         time.Time  `json:"due_date" db:"due_date"`
	PaidAt          *time.Time `json:"paid_at" db:"paid_at"`
	AttemptedAt     *time.Time `json:"attempted_at" db:"attempted_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type UsageRecord struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	SubscriptionID string    `json:"subscription_id" db:"subscription_id"`
	MeterName      string    `json:"meter_name" db:"meter_name"`
	Quantity       int64     `json:"quantity" db:"quantity"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	SentToStripe   bool      `json:"sent_to_stripe" db:"sent_to_stripe"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type PaymentMethod struct {
	ID                    string    `json:"id" db:"id"`
	UserID                string    `json:"user_id" db:"user_id"`
	StripePaymentMethodID string    `json:"stripe_payment_method_id" db:"stripe_payment_method_id"`
	Type                  string    `json:"type" db:"type"`
	Last4                 string    `json:"last4" db:"last4"`
	ExpiryMonth           int       `json:"expiry_month" db:"expiry_month"`
	ExpiryYear            int       `json:"expiry_year" db:"expiry_year"`
	IsDefault             bool      `json:"is_default" db:"is_default"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

type BillingEvent struct {
	ID            string                 `json:"id" db:"id"`
	UserID        string                 `json:"user_id" db:"user_id"`
	EventType     string                 `json:"event_type" db:"event_type"`
	StripeEventID string                 `json:"stripe_event_id" db:"stripe_event_id"`
	Data          map[string]interface{} `json:"data" db:"data"`
	Processed     bool                   `json:"processed" db:"processed"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	ProcessedAt   *time.Time             `json:"processed_at" db:"processed_at"`
}

// DTO structures for API responses
type SubscriptionResponse struct {
	Subscription *Subscription `json:"subscription"`
	Invoice      *Invoice      `json:"current_invoice,omitempty"`
	Usage        []UsageRecord `json:"usage,omitempty"`
}

type WebhookPayload struct {
	ID       string                 `json:"id"`
	Object   string                 `json:"object"`
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Created  int64                  `json:"created"`
	LiveMode bool                   `json:"livemode"`
}
