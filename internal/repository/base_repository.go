package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{
		db: db,
	}
}

func (r *BaseRepository) Executor(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TxKey).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *BaseRepository) GenerateSortableID(length int) string {
	// 1. ดึงเวลาปัจจุบันระดับนาโนวินาที (ประทับเวลา)
	// จะได้เลขประมาณ 15-16 หลักในรูป Hex ซึ่งเรียงตามเวลาเสมอ
	now := time.Now().UnixNano()
	timestampHex := fmt.Sprintf("%x", now) // เช่น "18f5d1a2b3c"

	// 2. คำนวณหาจำนวนตัวอักษรที่ต้องสุ่มเพิ่มเพื่อให้ครบ length
	randomLength := length - len(timestampHex)

	// 3. สุ่มตัวอักษรที่เหลือ (Secure Random)
	randomPart := r.secureRandomString(randomLength)

	return timestampHex + randomPart
}

func (r *BaseRepository) secureRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "" // Fallback
	}
	for i := 0; i < n; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
