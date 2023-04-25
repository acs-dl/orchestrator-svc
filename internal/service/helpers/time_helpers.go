package helpers

import "time"

func SumStringDurationWithDuration(first string, secondDuration time.Duration) (time.Duration, error) {
	firstDuration, err := time.ParseDuration(first)
	if err != nil {
		return time.Duration(0), err
	}

	return firstDuration + secondDuration, nil
}

func RoundDuration(durationString string) (time.Duration, error) {
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return time.Duration(0), err
	}

	return duration.Round(time.Second), nil
}
