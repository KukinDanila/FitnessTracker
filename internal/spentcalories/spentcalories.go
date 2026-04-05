package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
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

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("invalid data format: %s", data)
	}

	activity := s[1]

	duration, err := time.ParseDuration(s[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid data format: %s", data)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("invalid data format: %s", data)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	return height * stepLengthCoefficient * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
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
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	if errCalc != nil {
		log.Println(errCalc)
		return "", errCalc
	}

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), dist, speed, calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case steps <= 0:
		return 0.0, fmt.Errorf("количество шагов должно быть больше нуля")
	case weight <= 0:
		return 0.0, fmt.Errorf("вес должен быть больше нуля")
	case height <= 0:
		return 0.0, fmt.Errorf("рост должен быть больше нуля")
	case duration <= 0:
		return 0.0, fmt.Errorf("продолжительность должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	calories := (weight * speed * duration.Minutes()) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case steps <= 0:
		return 0.0, fmt.Errorf("количество шагов должно быть больше нуля")
	case weight <= 0:
		return 0.0, fmt.Errorf("вес должен быть больше нуля")
	case height <= 0:
		return 0.0, fmt.Errorf("рост должен быть больше нуля")
	case duration <= 0:
		return 0.0, fmt.Errorf("продолжительность должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	calories := (weight * speed * duration.Minutes()) / minInH * walkingCaloriesCoefficient
	return calories, nil
}
