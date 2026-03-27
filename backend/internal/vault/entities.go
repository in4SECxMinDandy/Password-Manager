package vault

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EntryType string

const (
	EntryTypeLogin    EntryType = "login"
	EntryTypeNote     EntryType = "note"
	EntryTypeCard     EntryType = "card"
	EntryTypeIdentity EntryType = "identity"
)

type Vault struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VaultEntry struct {
	ID        uuid.UUID `json:"id"`
	VaultID   uuid.UUID `json:"vault_id"`
	Type      EntryType `json:"type"`
	Title     string    `json:"title"`
	Data      []byte    `json:"data"`
	Favorite  bool      `json:"favorite"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginData struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URL      string `json:"url,omitempty"`
	Notes    string `json:"notes,omitempty"`
	TOTPSecret string `json:"totp_secret,omitempty"`
}

type NoteData struct {
	Content string `json:"content"`
}

type CardData struct {
	CardholderName string `json:"cardholder_name,omitempty"`
	Number         string `json:"number,omitempty"`
	ExpiryMonth    string `json:"expiry_month,omitempty"`
	ExpiryYear     string `json:"expiry_year,omitempty"`
	CVV            string `json:"cvv,omitempty"`
	Pin            string `json:"pin,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

type IdentityData struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Address1    string `json:"address_1,omitempty"`
	Address2    string `json:"address_2,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country,omitempty"`
}

func (e *VaultEntry) GetLoginData() (*LoginData, error) {
	if e.Type != EntryTypeLogin {
		return nil, nil
	}
	var data LoginData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (e *VaultEntry) GetNoteData() (*NoteData, error) {
	if e.Type != EntryTypeNote {
		return nil, nil
	}
	var data NoteData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (e *VaultEntry) GetCardData() (*CardData, error) {
	if e.Type != EntryTypeCard {
		return nil, nil
	}
	var data CardData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (e *VaultEntry) GetIdentityData() (*IdentityData, error) {
	if e.Type != EntryTypeIdentity {
		return nil, nil
	}
	var data IdentityData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

type CreateVaultRequest struct {
	Name string `json:"name"`
}

type UpdateVaultRequest struct {
	Name string `json:"name"`
}

type CreateEntryRequest struct {
	Type     EntryType `json:"type"`
	Title    string    `json:"title"`
	Data     string    `json:"data"`
	Favorite bool      `json:"favorite"`
}

type UpdateEntryRequest struct {
	Title    string    `json:"title"`
	Data     string    `json:"data"`
	Favorite *bool     `json:"favorite,omitempty"`
}

type VaultListResponse struct {
	Vaults []Vault `json:"vaults"`
	Total  int     `json:"total"`
}

type EntryListResponse struct {
	Entries []VaultEntry `json:"entries"`
	Total   int          `json:"total"`
}
