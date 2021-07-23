package domain

import "time"

type DateRange struct {
	From time.Time
	To   time.Time
}

func NewDateRange(from time.Time, to time.Time) DateRange {
	return DateRange{
		From: from,
		To:   to,
	}
}

func NewDateRangeFromString(fromString string, toString string) *DateRange {
	from, err := DateFrom(fromString)
	if err != nil {
		return nil
	}

	to, err := DateFrom(toString)
	if err != nil {
		return nil
	}

	dateRange := NewDateRange(from, to)
	return &dateRange
}

func NewDateRangeWithCap(from time.Time, to time.Time, cap time.Time) DateRange {
	dateTo := to
	if cap.Before(dateTo) {
		dateTo = cap
	}

	return DateRange{
		From: from,
		To:   dateTo,
	}
}

func (d *DateRange) In(t time.Time) bool {
	return d.From.Before(t) && d.To.After(t)
}
