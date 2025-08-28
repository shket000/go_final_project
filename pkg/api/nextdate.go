package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GET /api/nextdate?now=20060102&date=20060102&repeat=<rule>
//
// Если now не задан, берётся сегодняшняя дата (локаль не важна — только календарный день).
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	nowStr := q.Get("now")
	dstart := q.Get("date")
	repeat := q.Get("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.ParseInLocation(DateFmt, nowStr, time.Local)
		if err != nil {
			writeError(w, errors.New("invalid now date"), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(next))
}

// NextDate вычисляет ближайшую дату > now по правилу repeat, начиная от dstart.
// dstart и возвращаемая дата в формате 20060102.
// Поддержаны правила:
//   - ""         -> ошибка
//   - "d N"      -> каждые N дней, 1..400
//   - "y"        -> ежегодно
//   - "w list"   -> дни недели 1..7 (1=Пн, 7=Вс), через запятую
//   - "m days [months]" -> дни месяца (1..31, -1, -2), опционально месяцы (1..12) через пробел после списка дней
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if dstart == "" {
		return "", errors.New("empty dstart")
	}
	start, err := time.ParseInLocation(DateFmt, dstart, time.Local)
	if err != nil {
		return "", errors.New("invalid dstart format")
	}
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid d rule")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n <= 0 || n > 400 {
			return "", errors.New("invalid d interval")
		}
		date := start
		for {
			date = date.AddDate(0, 0, n)
			if date.After(stripTime(now)) {
				return date.Format(DateFmt), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid y rule")
		}
		date := start
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(stripTime(now)) {
				return date.Format(DateFmt), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid w rule")
		}
		daysList := strings.Split(parts[1], ",")
		var dow [8]bool // 1..7
		for _, s := range daysList {
			s = strings.TrimSpace(s)
			if s == "" {
				return "", errors.New("invalid w list")
			}
			v, err := strconv.Atoi(s)
			if err != nil || v < 1 || v > 7 {
				return "", errors.New("invalid w value")
			}
			dow[v] = true
		}
		date := start
		for {
			date = date.AddDate(0, 0, 1)
			// Go: Monday=1 ... Sunday=0; нам нужно 1..7 (Пн..Вс)
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			if dow[wd] && date.After(stripTime(now)) {
				return date.Format(DateFmt), nil
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid m rule")
		}
		// days
		dayTok := strings.Split(parts[1], ",")
		var mdays [32]bool // 1..31
		useLast := false
		usePrevLast := false
		for _, s := range dayTok {
			s = strings.TrimSpace(s)
			if s == "" {
				return "", errors.New("invalid m days")
			}
			if s == "-1" {
				useLast = true
				continue
			}
			if s == "-2" {
				usePrevLast = true
				continue
			}
			v, err := strconv.Atoi(s)
			if err != nil || v < 1 || v > 31 {
				return "", errors.New("invalid m day")
			}
			mdays[v] = true
		}
		// months (optional)
		var mfilter [13]bool // 1..12
		hasMonths := false
		if len(parts) == 3 {
			monTok := strings.Split(parts[2], ",")
			for _, s := range monTok {
				s = strings.TrimSpace(s)
				v, err := strconv.Atoi(s)
				if err != nil || v < 1 || v > 12 {
					return "", errors.New("invalid m month")
				}
				mfilter[v] = true
				hasMonths = true
			}
		}

		date := start
		for {
			date = date.AddDate(0, 0, 1)
			if date.After(stripTime(now)) {
				m := int(date.Month())
				if hasMonths && !mfilter[m] {
					continue
				}
				d := date.Day()
				ok := mdays[d]
				if !ok && useLast {
					if isLastDayOfMonth(date) {
						ok = true
					}
				}
				if !ok && usePrevLast {
					if isPrevLastDayOfMonth(date) {
						ok = true
					}
				}
				if ok {
					return date.Format(DateFmt), nil
				}
			}
		}

	default:
		return "", errors.New("unsupported repeat rule")
	}
}

func stripTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func isLastDayOfMonth(t time.Time) bool {
	y, m, _ := t.Date()
	firstNext := time.Date(y, m+1, 1, 0, 0, 0, 0, t.Location())
	return t.AddDate(0, 0, 1).Equal(firstNext)
}

func isPrevLastDayOfMonth(t time.Time) bool {
	// true, если t + 2 дня == первый день следующего месяца
	y, m, _ := t.Date()
	firstNext := time.Date(y, m+1, 1, 0, 0, 0, 0, t.Location())
	return t.AddDate(0, 0, 2).Equal(firstNext)
}
