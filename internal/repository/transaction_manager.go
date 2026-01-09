package repository

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey struct{}

var TxKey = ctxKey{}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (m *transactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	// ใช้ m.db.WithContext(ctx) เพื่อให้ Transaction รู้จัก Timeout/Cancellation จากภายนอก
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// สร้าง Context ใหม่โดยต่อยอดจาก ctx เดิม (เผื่อใน ctx มีค่าอื่นๆ เช่น UserID, TraceID)
		// แล้วฝาก tx เข้าไปใน Key ที่เรากำหนด
		txCtx := context.WithValue(ctx, TxKey, tx)

		// ส่ง txCtx ที่มีค่า tx ล่าสุดให้ Service ใช้งาน
		return fn(txCtx)
	})
}
