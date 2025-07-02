package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"pagemagic/billing-svc/internal/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/sub"
)

type BillingService struct {
	db *sql.DB
}

func NewBillingService(db *sql.DB) *BillingService {
	return &BillingService{db: db}
}

func (s *BillingService) CreateCustomer(ctx context.Context, userID, email string) (*models.User, error) {
	// Create Stripe customer
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Metadata: map[string]string{
			"user_id": userID,
		},
	}

	stripeCustomer, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Save to database
	query := `
		INSERT INTO users (id, email, stripe_customer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			stripe_customer_id = $3,
			updated_at = $5
		RETURNING id, email, stripe_customer_id, created_at, updated_at
	`

	now := time.Now()
	user := &models.User{}

	err = s.db.QueryRowContext(ctx, query, userID, email, stripeCustomer.ID, now, now).Scan(
		&user.ID, &user.Email, &user.StripeCustomerID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (s *BillingService) CreateSubscription(ctx context.Context, userID, planID string) (*models.Subscription, error) {
	// Get user's Stripe customer ID
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Create Stripe subscription
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(user.StripeCustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(planID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
		Expand: []*string{
			stripe.String("latest_invoice.payment_intent"),
		},
	}

	stripeSub, err := sub.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe subscription: %w", err)
	}

	// Save to database
	subscription := &models.Subscription{
		ID:                   fmt.Sprintf("sub_%d", time.Now().UnixNano()),
		UserID:               userID,
		StripeSubscriptionID: stripeSub.ID,
		Status:               string(stripeSub.Status),
		PlanID:               planID,
		CurrentPeriodStart:   time.Unix(stripeSub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(stripeSub.CurrentPeriodEnd, 0),
		CancelAtPeriodEnd:    stripeSub.CancelAtPeriodEnd,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	query := `
		INSERT INTO subscriptions (id, user_id, stripe_subscription_id, status, plan_id, 
			current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = s.db.ExecContext(ctx, query,
		subscription.ID, subscription.UserID, subscription.StripeSubscriptionID,
		subscription.Status, subscription.PlanID, subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd, subscription.CancelAtPeriodEnd,
		subscription.CreatedAt, subscription.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	return subscription, nil
}

func (s *BillingService) HandleWebhook(ctx context.Context, eventType string, data map[string]interface{}) error {
	switch eventType {
	case "invoice.payment_succeeded":
		return s.handleInvoicePaymentSucceeded(ctx, data)
	case "invoice.payment_failed":
		return s.handleInvoicePaymentFailed(ctx, data)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, data)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, data)
	default:
		log.Printf("Unhandled webhook event type: %s", eventType)
	}
	return nil
}

func (s *BillingService) handleInvoicePaymentSucceeded(ctx context.Context, data map[string]interface{}) error {
	// Extract invoice data
	invoiceData, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid invoice data")
	}

	stripeInvoiceID, _ := invoiceData["id"].(string)
	subscriptionID, _ := invoiceData["subscription"].(string)

	// Update invoice status
	query := `
		UPDATE invoices 
		SET status = 'paid', paid_at = $1, updated_at = $2
		WHERE stripe_invoice_id = $3
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, now, now, stripeInvoiceID)
	if err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	// Update subscription status if needed
	if subscriptionID != "" {
		query = `
			UPDATE subscriptions 
			SET status = 'active', updated_at = $1
			WHERE stripe_subscription_id = $2
		`
		_, err = s.db.ExecContext(ctx, query, now, subscriptionID)
		if err != nil {
			log.Printf("Failed to update subscription status: %v", err)
		}
	}

	return nil
}

func (s *BillingService) handleInvoicePaymentFailed(ctx context.Context, data map[string]interface{}) error {
	// Extract invoice data
	invoiceData, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid invoice data")
	}

	stripeInvoiceID, _ := invoiceData["id"].(string)

	// Update invoice status
	query := `
		UPDATE invoices 
		SET status = 'payment_failed', attempted_at = $1, updated_at = $2
		WHERE stripe_invoice_id = $3
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, now, now, stripeInvoiceID)
	if err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	// TODO: Implement dunning logic - suspend services after failed payments

	return nil
}

func (s *BillingService) handleSubscriptionUpdated(ctx context.Context, data map[string]interface{}) error {
	subData, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid subscription data")
	}

	stripeSubID, _ := subData["id"].(string)
	status, _ := subData["status"].(string)
	cancelAtPeriodEnd, _ := subData["cancel_at_period_end"].(bool)

	query := `
		UPDATE subscriptions 
		SET status = $1, cancel_at_period_end = $2, updated_at = $3
		WHERE stripe_subscription_id = $4
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, status, cancelAtPeriodEnd, now, stripeSubID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (s *BillingService) handleSubscriptionDeleted(ctx context.Context, data map[string]interface{}) error {
	subData, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid subscription data")
	}

	stripeSubID, _ := subData["id"].(string)

	query := `
		UPDATE subscriptions 
		SET status = 'canceled', updated_at = $1
		WHERE stripe_subscription_id = $2
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, now, stripeSubID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (s *BillingService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT id, email, stripe_customer_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.StripeCustomerID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *BillingService) GetUserSubscriptions(ctx context.Context, userID string) ([]models.Subscription, error) {
	query := `
		SELECT id, user_id, stripe_subscription_id, status, plan_id,
			current_period_start, current_period_end, cancel_at_period_end,
			created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.StripeSubscriptionID, &sub.Status, &sub.PlanID,
			&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd,
			&sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (s *BillingService) SaveBillingEvent(ctx context.Context, stripeEventID, eventType string, data map[string]interface{}) error {
	// Extract user ID from the event data if possible
	userID := s.extractUserIDFromEvent(data)

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	query := `
		INSERT INTO billing_events (id, user_id, event_type, stripe_event_id, data, processed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (stripe_event_id) DO NOTHING
	`

	eventID := fmt.Sprintf("evt_%d", time.Now().UnixNano())
	now := time.Now()

	_, err = s.db.ExecContext(ctx, query, eventID, userID, eventType, stripeEventID, dataJSON, false, now)
	if err != nil {
		return fmt.Errorf("failed to save billing event: %w", err)
	}

	return nil
}

func (s *BillingService) extractUserIDFromEvent(data map[string]interface{}) string {
	// Try to extract user ID from various possible locations in the event data
	if obj, ok := data["object"].(map[string]interface{}); ok {
		if metadata, ok := obj["metadata"].(map[string]interface{}); ok {
			if userID, ok := metadata["user_id"].(string); ok {
				return userID
			}
		}
	}
	return ""
}
