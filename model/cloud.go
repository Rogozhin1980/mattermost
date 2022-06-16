// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"os"
	"strconv"
	"strings"
)

const (
	EventTypeFailedPayment                = "failed-payment"
	EventTypeFailedPaymentNoCard          = "failed-payment-no-card"
	EventTypeSendAdminWelcomeEmail        = "send-admin-welcome-email"
	EventTypeSendUpgradeConfirmationEmail = "send-upgrade-confirmation-email"
	EventTypeSubscriptionChanged          = "subscription-changed"
	EventTypeTrialWillEnd                 = "trial-will-end"
	EventTypeTrialEnded                   = "trial-ended"
)

var MockCWS string

type BillingScheme string

const (
	BillingSchemePerSeat    = BillingScheme("per_seat")
	BillingSchemeFlatFee    = BillingScheme("flat_fee")
	BillingSchemeSalesServe = BillingScheme("sales_serve")
)

type RecurringInterval string

const (
	RecurringIntervalYearly  = RecurringInterval("year")
	RecurringIntervalMonthly = RecurringInterval("month")
)

type SubscriptionFamily string

const (
	SubscriptionFamilyCloud  = SubscriptionFamily("cloud")
	SubscriptionFamilyOnPrem = SubscriptionFamily("on-prem")
)

const defaultCloudNotifyAdminCoolOffDays = 30

// Product model represents a product on the cloud system.
type Product struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	PricePerSeat      float64            `json:"price_per_seat"`
	AddOns            []*AddOn           `json:"add_ons"`
	SKU               string             `json:"sku"`
	PriceID           string             `json:"price_id"`
	Family            SubscriptionFamily `json:"product_family"`
	RecurringInterval RecurringInterval  `json:"recurring_interval"`
	BillingScheme     BillingScheme      `json:"billing_scheme"`
}

// AddOn represents an addon to a product.
type AddOn struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	DisplayName  string  `json:"display_name"`
	PricePerSeat float64 `json:"price_per_seat"`
}

// StripeSetupIntent represents the SetupIntent model from Stripe for updating payment methods.
type StripeSetupIntent struct {
	ID           string `json:"id"`
	ClientSecret string `json:"client_secret"`
}

// ConfirmPaymentMethodRequest contains the fields for the customer payment update API.
type ConfirmPaymentMethodRequest struct {
	StripeSetupIntentID string `json:"stripe_setup_intent_id"`
	SubscriptionID      string `json:"subscription_id"`
}

// Customer model represents a customer on the system.
type CloudCustomer struct {
	CloudCustomerInfo
	ID             string         `json:"id"`
	CreatorID      string         `json:"creator_id"`
	CreateAt       int64          `json:"create_at"`
	BillingAddress *Address       `json:"billing_address"`
	CompanyAddress *Address       `json:"company_address"`
	PaymentMethod  *PaymentMethod `json:"payment_method"`
}

type ValidateBusinessEmailRequest struct {
	Email string `json:"email"`
}

// CloudCustomerInfo represents editable info of a customer.
type CloudCustomerInfo struct {
	Name             string `json:"name"`
	Email            string `json:"email,omitempty"`
	ContactFirstName string `json:"contact_first_name,omitempty"`
	ContactLastName  string `json:"contact_last_name,omitempty"`
	NumEmployees     int    `json:"num_employees"`
}

// Address model represents a customer's address.
type Address struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

// PaymentMethod represents methods of payment for a customer.
type PaymentMethod struct {
	Type      string `json:"type"`
	LastFour  string `json:"last_four"`
	ExpMonth  int    `json:"exp_month"`
	ExpYear   int    `json:"exp_year"`
	CardBrand string `json:"card_brand"`
	Name      string `json:"name"`
}

// Subscription model represents a subscription on the system.
type Subscription struct {
	ID          string   `json:"id"`
	CustomerID  string   `json:"customer_id"`
	ProductID   string   `json:"product_id"`
	AddOns      []string `json:"add_ons"`
	StartAt     int64    `json:"start_at"`
	EndAt       int64    `json:"end_at"`
	CreateAt    int64    `json:"create_at"`
	Seats       int      `json:"seats"`
	Status      string   `json:"status"`
	DNS         string   `json:"dns"`
	IsPaidTier  string   `json:"is_paid_tier"`
	LastInvoice *Invoice `json:"last_invoice"`
	IsFreeTrial string   `json:"is_free_trial"`
	TrialEndAt  int64    `json:"trial_end_at"`
}

// GetWorkSpaceNameFromDNS returns the work space name. For example from test.mattermost.cloud.com, it returns test
func (s *Subscription) GetWorkSpaceNameFromDNS() string {
	return strings.Split(s.DNS, ".")[0]
}

// Invoice model represents a cloud invoice
type Invoice struct {
	ID                 string             `json:"id"`
	Number             string             `json:"number"`
	CreateAt           int64              `json:"create_at"`
	Total              int64              `json:"total"`
	Tax                int64              `json:"tax"`
	Status             string             `json:"status"`
	Description        string             `json:"description"`
	PeriodStart        int64              `json:"period_start"`
	PeriodEnd          int64              `json:"period_end"`
	SubscriptionID     string             `json:"subscription_id"`
	Items              []*InvoiceLineItem `json:"line_items"`
	CurrentProductName string             `json:"current_product_name"`
}

// InvoiceLineItem model represents a cloud invoice lineitem tied to an invoice.
type InvoiceLineItem struct {
	PriceID      string                 `json:"price_id"`
	Total        int64                  `json:"total"`
	Quantity     float64                `json:"quantity"`
	PricePerUnit int64                  `json:"price_per_unit"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type CWSWebhookPayload struct {
	Event                             string               `json:"event"`
	FailedPayment                     *FailedPayment       `json:"failed_payment"`
	CloudWorkspaceOwner               *CloudWorkspaceOwner `json:"cloud_workspace_owner"`
	ProductLimits                     *ProductLimits       `json:"product_limits"`
	Subscription                      *Subscription        `json:"subscription"`
	SubscriptionTrialEndUnixTimeStamp int64                `json:"trial_end_time_stamp"`
}

type FailedPayment struct {
	CardBrand      string `json:"card_brand"`
	LastFour       string `json:"last_four"`
	FailureMessage string `json:"failure_message"`
}

// CloudWorkspaceOwner is part of the CWS Webhook payload that contains information about the user that created the workspace from the CWS
type CloudWorkspaceOwner struct {
	UserName string `json:"username"`
}
type SubscriptionChange struct {
	ProductID string `json:"product_id"`
}

type BoardsLimits struct {
	Cards *int `json:"cards"`
	Views *int `json:"views"`
}

type FilesLimits struct {
	TotalStorage *int64 `json:"total_storage"`
}

type IntegrationsLimits struct {
	Enabled *int `json:"enabled"`
}

type MessagesLimits struct {
	History *int `json:"history"`
}

type TeamsLimits struct {
	Active *int `json:"active"`
}

type ProductLimits struct {
	Boards       *BoardsLimits       `json:"boards,omitempty"`
	Files        *FilesLimits        `json:"files,omitempty"`
	Integrations *IntegrationsLimits `json:"integrations,omitempty"`
	Messages     *MessagesLimits     `json:"messages,omitempty"`
	Teams        *TeamsLimits        `json:"teams,omitempty"`
}

type NotifyAdminToUpgradeRequest struct {
	CurrentTeamId string `json:"current_team_id"`
}

type UserInfo struct {
	UserID    string
	Timestamp int64
}

type AlreadyCloudNotifiedAdminUsersInfo struct {
	Info []UserInfo
}

func (a *AlreadyCloudNotifiedAdminUsersInfo) CanNotify(ID string) bool {
	coolOffPeriodDaysEnv := os.Getenv("MM_CLOUD_NOTIFY_ADMIN_COOL_OFF_DAYS")
	coolOffPeriodDays, parseError := strconv.ParseFloat(coolOffPeriodDaysEnv, 64)
	if parseError != nil {
		coolOffPeriodDays = defaultCloudNotifyAdminCoolOffDays
	}
	daysToMillis := coolOffPeriodDays * 24 * 60 * 60 * 1000
	for _, i := range a.Info {
		if i.UserID == ID {
			timeDiff := GetMillis() - i.Timestamp
			if timeDiff >= int64(daysToMillis) {
				return true
			}
			return false
		}
	}

	return true
}

func (a *AlreadyCloudNotifiedAdminUsersInfo) Upsert(ID string) []UserInfo {
	for ind, i := range a.Info {
		if i.UserID == ID {
			currentUserInfo := a.Info[ind]
			currentUserInfo.Timestamp = GetMillis()
			a.Info[ind] = currentUserInfo
			return a.Info
		}
	}

	a.Info = append(a.Info, UserInfo{
		UserID:    ID,
		Timestamp: GetMillis(),
	})

	return a.Info
}
