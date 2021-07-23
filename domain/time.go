package domain

import "time"

func UTC(t time.Time) time.Time {
	return t.In(time.UTC)
}

func UTCFrom(ts string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", ts)
	if err != nil {
		return time.Time{}, err
	}

	return jst2utc(t), nil
}

func DateFrom(ts string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func JSTStringFromDateTime(now time.Time) string {
	return utc2jst(now).Format("2006-01-02 15:04:05")
}

func StringFromDate(now time.Time) string {
	return utc2jst(now).Format("2006-01-02")
}

func NowUTC() time.Time {
	return time.Now().In(time.UTC)
}

func utc2jst(t time.Time) time.Time {
	return t.In(time.FixedZone("Asia/Tokyo", 9*60*60))
}

func jst2utc(t time.Time) time.Time {
	return t.In(time.UTC).Add(-9 * 60 * 60 * time.Second)
}
