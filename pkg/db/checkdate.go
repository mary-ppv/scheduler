package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat field")
	}

	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("incorrect start date")
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	if start.Month() == time.February && start.Day() == 29 && rule == "y" {
		date := time.Date(start.Year()+1, 3, 1, 0, 0, 0, 0, start.Location())
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(DateFormat), nil
	}

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("incorrect format of rule d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("interval should be between 1 and 400")
		}

		date := start.AddDate(0, 0, interval)
		for !date.After(now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("incorrect format of rule d y")
		}

		date := start.AddDate(1, 0, 0)
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(DateFormat), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("incorrect format of rule d w")
		}
		dayStrs := strings.Split(parts[1], ",")
		days := make([]int, 0, len(dayStrs))
		for _, s := range dayStrs {
			day, err := strconv.Atoi(s)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("incorrect day of the week")
			}
			days = append(days, day)
		}

		weekdayMap := map[int]time.Weekday{
			1: time.Monday,
			2: time.Tuesday,
			3: time.Wednesday,
			4: time.Thursday,
			5: time.Friday,
			6: time.Saturday,
			7: time.Sunday,
		}

		validDays := make(map[time.Weekday]bool)
		for _, d := range days {
			validDays[weekdayMap[d]] = true
		}

		date := start.AddDate(0, 0, 1)
		for !date.After(now) || !validDays[date.Weekday()] {
			date = date.AddDate(0, 0, 1)
			if date.Year()-now.Year() > 10 {
				return "", errors.New("can not find date for w")
			}
		}
		return date.Format(DateFormat), nil

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("incorrect format for rule m")
		}

		dayStrs := strings.Split(parts[1], ",")
		monthDays := make([]int, 0, len(dayStrs))
		for _, s := range dayStrs {
			day, err := strconv.Atoi(s)
			if err != nil || day < -2 || day == 0 || day > 31 {
				return "", errors.New("incorect day of month")
			}
			monthDays = append(monthDays, day)
		}

		var validMonths [13]bool
		if len(parts) == 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, s := range monthStrs {
				month, err := strconv.Atoi(s)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("incorrect month")
				}
				validMonths[month] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				validMonths[i] = true
			}
		}

		date := start.AddDate(0, 0, 1)
		for {
			if date.After(now) {
				year, month, day := date.Date()
				monthInt := int(month)
				lastDay := lastDayOfMonth(year, monthInt)

				matched := false
				for _, md := range monthDays {
					var targetDay int
					switch md {
					case -1:
						targetDay = lastDay
					case -2:
						targetDay = lastDay - 1
					default:
						targetDay = md
					}
					if day == targetDay {
						matched = true
						break
					}
				}

				if matched && validMonths[monthInt] {
					return date.Format(DateFormat), nil
				}
			}
			date = date.AddDate(0, 0, 1)
			if date.Year()-now.Year() > 10 {
				return "", errors.New("can not find date for m")
			}
		}

	default:
		return "", fmt.Errorf("unsupported rule: %s", rule)
	}
}

func lastDayOfMonth(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Invalid 'now' date: %v", err)})
			return
		}
	}

	if dateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "date parameter is required"})
		return
	}

	if repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "repeat parameter is required"})
		return
	}

	nextDateStr, err := NextDate(now, dateStr, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDateStr + "\n"))
}

func CheckDate(task *Task, now time.Time) error {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" || task.Date == "today" {
		task.Date = now.Format("20060102")
		return nil
	}

	if len(task.Date) != 8 {
		return fmt.Errorf("invalid date format: must be YYYYMMDD")
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	if date.Format("20060102") != task.Date {
		return fmt.Errorf("invalid date: %s", task.Date)
	}

	if !date.After(now) && task.Repeat == "" {
		task.Date = now.Format("20060102")
	}

	if !date.After(now) && task.Repeat != "" {
		nextDateStr, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("error computing next date: %v", err)
		}
		task.Date = nextDateStr
	}

	return nil
}

func AfterNow(t time.Time, now time.Time) bool {
	return t.After(now) || t.Equal(now)
}
