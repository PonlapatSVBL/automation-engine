package utils

import (
	"fmt"
	"strings"
	"time"
)

func CalculateDailyNextRun(
	now time.Time,
	runTime time.Time,
	loc *time.Location,
) (time.Time, error) {

	h, m, s := runTime.Hour(), runTime.Minute(), runTime.Second()

	candidate := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		h, m, s,
		0,
		loc,
	)

	if candidate.After(now) {
		return candidate, nil
	}

	return candidate.AddDate(0, 0, 1), nil
}

func CalculateWeeklyNextRun(
	now time.Time,
	runTime time.Time,
	dayOfWeek string,
	loc *time.Location,
) (time.Time, error) {

	var weekdayMap = map[string]time.Weekday{
		"sun": time.Sunday,
		"mon": time.Monday,
		"tue": time.Tuesday,
		"wed": time.Wednesday,
		"thu": time.Thursday,
		"fri": time.Friday,
		"sat": time.Saturday,
	}

	targetWeekday, ok := weekdayMap[strings.ToLower(dayOfWeek)]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid day_of_week: %s", dayOfWeek)
	}

	h, m, s := runTime.Hour(), runTime.Minute(), runTime.Second()

	// วันนี้ใน timezone ที่ถูกต้อง
	now = now.In(loc)

	// สร้าง candidate ของ "สัปดาห์นี้"
	candidate := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		h, m, s,
		0,
		loc,
	)

	// คำนวณจำนวนวันต้องขยับ
	daysDiff := int(targetWeekday - now.Weekday())
	if daysDiff < 0 {
		daysDiff += 7
	}

	candidate = candidate.AddDate(0, 0, daysDiff)

	// ถ้าวันตรง แต่เวลาเลยแล้ว → ขยับไปอาทิตย์หน้า
	if daysDiff == 0 && !candidate.After(now) {
		candidate = candidate.AddDate(0, 0, 7)
	}

	return candidate, nil
}

func CalculateMonthlyNextRun(
	now time.Time,
	runTime time.Time,
	dayOfMonth int,
	loc *time.Location,
) (time.Time, error) {

	if dayOfMonth < 1 || dayOfMonth > 31 {
		return time.Time{}, fmt.Errorf("invalid day_of_month: %d", dayOfMonth)
	}

	h, m, s := runTime.Hour(), runTime.Minute(), runTime.Second()
	now = now.In(loc)

	year, month := now.Year(), now.Month()

	for i := 0; i < 24; i++ { // max 2 ปี ป้องกัน infinite loop
		daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

		if dayOfMonth <= daysInMonth {
			candidate := time.Date(
				year,
				month,
				dayOfMonth,
				h, m, s,
				0,
				loc,
			)

			if candidate.After(now) {
				return candidate, nil
			}
		}

		// ขยับไปเดือนถัดไป
		month++
		if month > 12 {
			month = 1
			year++
		}
	}

	return time.Time{}, fmt.Errorf("unable to calculate next monthly run")
}

// internal/utils/date.go

func CalculateYearlyNextRun(
	now time.Time,
	runTime time.Time,
	dayOfMonth int,
	monthOfYear int,
	loc *time.Location,
) (time.Time, error) {
	// ตรวจสอบความถูกต้องของเดือนและวันเบื้องต้น
	if monthOfYear < 1 || monthOfYear > 12 {
		return time.Time{}, fmt.Errorf("invalid month_of_year: %d", monthOfYear)
	}

	h, m, s := runTime.Hour(), runTime.Minute(), runTime.Second()
	now = now.In(loc)

	year := now.Year()

	// พยายามสร้าง candidate ในปีปัจจุบันก่อน
	// ใช้ monthOfYear และ dayOfMonth ที่ได้รับมาจากพารามิเตอร์
	candidate := time.Date(
		year,
		time.Month(monthOfYear),
		dayOfMonth,
		h, m, s,
		0,
		loc,
	)

	// ถ้า candidate ที่สร้างขึ้น "ไม่มากกว่า" เวลาปัจจุบัน (คือผ่านมาแล้วในรอบปีนี้)
	// ให้ขยับไปปีหน้า
	if !candidate.After(now) {
		candidate = candidate.AddDate(1, 0, 0)
	}

	return candidate, nil
}
