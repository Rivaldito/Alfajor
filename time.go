package alfajor

import (
	"strconv"
	"time"
)

func getDateTime() string {
	year, month, day := time.Now().Date()
	dateTime := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	return dateTime

}

func getHourAndMinutes() string {
	hours, minutes, seconds := time.Now().Clock()
	time := strconv.Itoa(hours) + ":" + strconv.Itoa(int(minutes)) + ":" + strconv.Itoa(seconds)
	return time
}
