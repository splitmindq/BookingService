package entity

import (
	"github.com/Rhymond/go-money"
	"time"
)

type PaymentMethod string

const (
	PaymentMethodCard PaymentMethod = "card"
	PaymentMethodBank PaymentMethod = "bank"
)

type PaymentStatus string

const (
	PendingPaymentStatus PaymentStatus = "pending"
	SuccessPaymentStatus PaymentStatus = "success"
	FailedPaymentStatus  PaymentStatus = "failed"
)

type Payment struct {
	BookingId int
	Amount    *money.Money
	// sql --> amount & currency
	//type Money struct {
	//    amount   Amount    `db:"amount"`
	//    currency *Currency `db:"currency"`
	//}
	paymentMethod PaymentMethod
	paymentStatus PaymentStatus
	transactionId int64
	CreatedAt     time.Time
}
