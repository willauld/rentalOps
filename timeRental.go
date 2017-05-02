package main

import "time"

// RentalNowTime is a Testing replacement for time.Now(). It is set in initTimeRental()
var RentalNowTime *time.Time

func initTimeRental(nowTime time.Time) {
	RentalNowTime = new(time.Time)
	*RentalNowTime = nowTime
}
func timeNowRental() time.Time {
	if RentalNowTime != nil {
		return *RentalNowTime
	}
	return time.Now()
}
