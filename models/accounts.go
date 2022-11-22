package models

import "github.com/google/uuid"

type AccountData struct {
	ID             uuid.UUID          `json:"id,omitempty"`
	OrganisationID uuid.UUID          `json:"organisation_id,omitempty"`
	Type           AccountType        `json:"type"`
	Attributes     *AccountAttributes `json:"attributes"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountType string

const (
	Accounts AccountType = "accounts"
)

type AccountAttributes struct {
	AccountClassification   *AccountClassification `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool                  `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string                 `json:"account_number,omitempty"`
	AlternativeNames        []string               `json:"alternative_names,omitempty"`
	BankID                  string                 `json:"bank_id,omitempty"`
	BankIDCode              string                 `json:"bank_id_code,omitempty"`
	BaseCurrency            string                 `json:"base_currency,omitempty"`
	Bic                     string                 `json:"bic,omitempty"`
	Country                 *string                `json:"country,omitempty"`
	Iban                    string                 `json:"iban,omitempty"`
	JointAccount            *bool                  `json:"joint_account,omitempty"`
	Name                    []string               `json:"name,omitempty"`
	SecondaryIdentification string                 `json:"secondary_identification,omitempty"`
	Status                  *AccountStatus         `json:"status,omitempty"`
	Switched                *bool                  `json:"switched,omitempty"`
}

type AccountClassification string

const (
	Business AccountClassification = "Business"
	Personal                       = "Personal"
)

type AccountStatus string

const (
	Closed    AccountStatus = "closed"
	Confirmed               = "confirmed"
	Pending                 = "pending"
)
