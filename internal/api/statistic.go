package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"time"
	"trackly-backend/internal/models"
	"trackly-backend/internal/repositories"
)

type StatisticApi struct {
	habitRepo     *repositories.HabitRepository
	statisticRepo *repositories.StatisticRepository
}

func NewStatisticApi(habitRepo *repositories.HabitRepository, statisticRepo *repositories.StatisticRepository) *StatisticApi {
	return &StatisticApi{habitRepo: habitRepo, statisticRepo: statisticRepo}
}

func (s StatisticApi) GetApiHabitsHabitIdStatistic(ctx echo.Context, habitId int, params GetApiHabitsHabitIdStatisticParams) error {
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{})
	}
	dateFrom := params.DateFrom.Time
	dateTo := params.DateTo.Time.Add(time.Hour*24 - 1)

	userId := ctx.Get("user_id").(int)

	habit, err := s.habitRepo.GetHabitWithStatInInterval(habitId, userId, dateFrom, dateTo)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
	}

	// сортируем
	sort.Slice(habit.HabitScore, func(i, j int) bool {
		return habit.HabitScore[i].DateTime.Before(habit.HabitScore[i].DateTime)
	})
	currentPlan := findCurrentPlan(habit.Plans)

	switch params.GroupBy {

	case Day:
		byDay := periodsByDay(habit, dateFrom, dateTo)
		return ctx.JSON(http.StatusOK, HabitStatisticResponse{
			GroupBy:  params.GroupBy,
			Period:   byDay,
			PlanUnit: PlanUnit(fmt.Sprintf("%v", currentPlan.PlanUnit)),
		})
	case Month:
		return ctx.JSON(http.StatusOK, HabitStatisticResponse{
			GroupBy:  params.GroupBy,
			Period:   periodsByMonth(habit, params.DateFrom.Time, params.DateTo.Time),
			PlanUnit: PlanUnit(fmt.Sprintf("%v", currentPlan.PlanUnit)),
		})
	case Year:
		return ctx.JSON(http.StatusOK, HabitStatisticResponse{
			GroupBy:  params.GroupBy,
			Period:   periodsByYear(habit, params.DateFrom.Time, params.DateTo.Time),
			PlanUnit: PlanUnit(fmt.Sprintf("%v", currentPlan.PlanUnit)),
		})
	}

	return ctx.JSON(http.StatusOK, nil)

}

func periodsByDay(habit *models.Habit, from time.Time, to time.Time) []PeriodValue {
	// Создаем карту для суммирования значений по дням
	dailySums := make(map[string]int)
	for _, score := range habit.HabitScore {
		dayKey := score.DateTime.Format(time.DateOnly) // Формат: ГГГГ-ММ-ДД
		dailySums[dayKey] += score.Value
	}

	// Находим минимальную и максимальную дату
	minDate, maxDate := findMinMaxDate(from, to)

	// Создаем итоговый массив
	var scores []PeriodValue
	for currentDate := minDate; !currentDate.After(maxDate); currentDate = currentDate.AddDate(0, 0, 1) {
		dayKey := currentDate.Format(time.DateOnly) // Формат: ГГГГ-ММ-ДД
		value := dailySums[dayKey]                  // Если день отсутствует в карте, значение будет 0
		scores = append(scores, PeriodValue{
			Interval: dayKey,
			Value:    value,
		})
	}
	return scores
}

func periodsByMonth(habit *models.Habit, from time.Time, to time.Time) []PeriodValue {
	monthlySums := make(map[string]int)
	for _, score := range habit.HabitScore {
		monthKey := score.DateTime.Format("2006-01") // Формат: ГГГГ-ММ
		monthlySums[monthKey] += score.Value
	}

	// Находим минимальную и максимальную дату
	minDate, maxDate := findMinMaxMonthDate(from, to)

	// Создаем итоговый массив
	scores := []PeriodValue{}
	for currentDate := minDate; !currentDate.After(maxDate); currentDate = currentDate.AddDate(0, 1, 0) {
		monthKey := currentDate.Format("2006-01") // Формат: ГГГГ-ММ
		value := monthlySums[monthKey]            // Если месяц отсутствует в карте, значение будет 0
		scores = append(scores, PeriodValue{
			Interval: monthKey,
			Value:    value,
		})
	}
	return scores
}

func periodsByYear(habit *models.Habit, from time.Time, to time.Time) []PeriodValue {
	// Создаем карту для суммирования значений по годам
	yearlySums := make(map[int]int)
	for _, score := range habit.HabitScore {
		year := score.DateTime.Year()
		yearlySums[year] += score.Value
	}

	// Находим минимальный и максимальный год
	minYear, maxYear := from.Year(), to.Year()

	// Создаем итоговый массив
	var scores []PeriodValue
	for year := minYear; year <= maxYear; year++ {
		value := yearlySums[year] // Если год отсутствует в карте, значение будет 0
		scores = append(scores, PeriodValue{
			Interval: fmt.Sprintf("%d", year),
			Value:    value,
		})
	}
	return scores
}

func (s StatisticApi) GetApiHabitsHabitIdStatisticTotal(ctx echo.Context, habitId int) error {

	userId := ctx.Get("user_id").(int)

	habit, err := s.habitRepo.GetHabitById(habitId, userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	uniqueDates := make(map[time.Time]bool)
	total := 0

	for _, score := range habit.HabitScore {
		total += score.Value
		uniqueDates[score.DateTime.UTC().Truncate(24*time.Hour)] = true
	}

	days := len(uniqueDates)
	average := 0
	if days > 0 {
		average = total / days
	}

	plansArr := habit.Plans

	var plan models.Plan
	if len(plansArr) > 0 {

		sort.Slice(plansArr, func(i, j int) bool {
			return plansArr[i].ID > plansArr[j].ID
		})
		plan = plansArr[0]
	}

	response := HabitStatisticTotalResponse{
		AveragePerDay: &average,
		PlanUnit:      (*PlanUnit)(plan.PlanUnit),
		Total:         &total,
	}

	return ctx.JSON(200, response)
}

func findCurrentPlan(plans []models.Plan) *models.Plan {
	var currentPlan *models.Plan
	for _, plan := range plans {
		if currentPlan == nil || plan.ID >= currentPlan.ID {
			currentPlan = &plan
		}
	}
	return currentPlan
}

func findMinMaxMonthDate(minDate, maxDate time.Time) (time.Time, time.Time) {
	// Приводим минимальную и максимальную дату к началу месяца
	minDate = time.Date(minDate.Year(), minDate.Month(), 1, 0, 0, 0, 0, minDate.Location())
	maxDate = time.Date(maxDate.Year(), maxDate.Month(), 1, 0, 0, 0, 0, maxDate.Location())
	return minDate, maxDate
}

func findMinMaxDate(minDate, maxDate time.Time) (time.Time, time.Time) {

	// Приводим минимальную и максимальную дату к началу дня
	minDate = time.Date(minDate.Year(), minDate.Month(), minDate.Day(), 0, 0, 0, 0, minDate.Location())
	maxDate = time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 0, 0, 0, 0, maxDate.Location())

	return minDate, maxDate
}
