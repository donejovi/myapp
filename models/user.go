package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `gorm:"unique" json:"phone_number"`
	Address     string    `json:"address"`
	PIN         string    `json:"pin"`
	Balance     float64   `json:"balance"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"update_date"`

	TopUps []TopUp `gorm:"foreignKey:UserID" json:"-"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}

type TopUp struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"top_up_id"`
	UserID        uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"-"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedDate   time.Time `json:"created_date"`
}

func (topUp *TopUp) BeforeCreate(tx *gorm.DB) (err error) {
	topUp.ID = uuid.New()
	return nil
}

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"payment_id"`
	UserID        uuid.UUID `gorm:"type:uuid;index" json:"-"`
	User          User      `gorm:"foreignKey:UserID" json:"-"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedDate   time.Time `json:"created_date"`
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	payment.ID = uuid.New()
	return nil
}

type Transfer struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"transfer_id"`
	FromUserID    uuid.UUID `gorm:"type:uuid;index" json:"-"`
	FromUser      User      `gorm:"foreignKey:FromUserID" json:"-"`
	ToUserID      uuid.UUID `gorm:"type:uuid;index" json:"-"`
	ToUser        User      `gorm:"foreignKey:ToUserID" json:"-"`
	Amount        float64   `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedDate   time.Time `json:"created_date"`
}

func (transfer *Transfer) BeforeCreate(tx *gorm.DB) (err error) {
	transfer.ID = uuid.New()
	return nil
}

type Transaction struct {
	UserID          uuid.UUID `gorm:"type:uuid;index" json:"-"`
	User            User      `gorm:"foreignKey:UserID" json:"-"`
	TransactionType string    `json:"transaction_type"` // CREDIT or DEBIT
	Amount          float64   `json:"amount"`
	Remarks         string    `json:"remarks"`
	BalanceBefore   float64   `json:"balance_before"`
	BalanceAfter    float64   `json:"balance_after"`
	CreatedDate     time.Time `json:"created_date"`
}
