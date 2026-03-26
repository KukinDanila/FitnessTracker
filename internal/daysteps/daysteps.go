package daysteps

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"time"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	s := strings.Split(data, ",")

	if len(s) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: %s", data)
	}

	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, 0, fmt.Errorf("неверное количество шагов: %s", s[0])
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля: %d", steps)
	}

	duration, err := time.ParseDuration(s[1])
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат продолжительности: %s", s[1])
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше нуля: %v", duration)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceKm := float64(steps) * stepLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)

	return result
}
