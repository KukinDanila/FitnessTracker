package spentcalories

import (
	"time"
	"strings"
	"log"
	"strconv"
	"fmt"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	s := strings.Split(data, ",")

	if len(s) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: %s", data)
	}

	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid data format: %s", data)
	}

	activity := s[1]

	duration, err := time.ParseDuration(s[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid data format: %s", data)
	}

	return steps, activity, duration, nil

}

func distance(steps int, height float64) float64 {
	lenSteps := (height * stepLengthCoefficient)

	distanceMetr := lenSteps * float64(steps)

	distanceKm := distanceMetr / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distKm := distance(steps, height)

	hours := duration.Hours()

	speed := distKm / float64(hours)

	return speed

}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	var errCalc error

	switch activity {
	case "Ходьба":
		calories, errCalc = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, errCalc = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("Неизвестный тип тренировки: %s", activity)
	}

	if errCalc != nil {
		log.Println(errCalc)
		return "", errCalc

	}

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	switch {
	case steps <= 0:
		return 0.0, nil
	case weight <= 0:
		return 0.0, nil
	case height <= 0:
		return 0.0, nil
	case duration <= 0:
		return 0.0, nil
	}

	speed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case steps <= 0:
		return 0.0, nil
	case weight <= 0:
		return 0.0, nil
	case height <= 0:
		return 0.0, nil
	case duration <= 0:
		return 0.0, nil
	}

	speed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	calories *= walkingCaloriesCoefficient

	return calories, nil


}
