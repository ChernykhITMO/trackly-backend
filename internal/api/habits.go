package api

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
	"sort"
	"time"
	"trackly-backend/internal/models"
	"trackly-backend/internal/repositories"
)

type HabitsApi struct {
	habitRepo     *repositories.HabitRepository
	planRepo      *repositories.PlanRepository
	statisticRepo *repositories.StatisticRepository
}

func NewHabitsApi(habitRepo *repositories.HabitRepository, planRepo *repositories.PlanRepository) *HabitsApi {
	return &HabitsApi{habitRepo: habitRepo, planRepo: planRepo}
}

func (h *HabitsApi) GetApiHabits(ctx echo.Context) error {

	userId := ctx.Get("user_id").(int)

	habits, err := h.habitRepo.GetHabitsByUserId(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	var response []Habit

	for _, habit := range habits {

		plansArr := habit.Plans
		var plan models.Plan
		if len(plansArr) > 0 {

			sort.Slice(plansArr, func(i, j int) bool {
				return plansArr[i].ID > plansArr[j].ID
			})
			plan = plansArr[0]
		}
		planUnit := plan.PlanUnit
		respPlan := Plan{
			Goal:     plan.Goal,
			PlanUnit: (*PlanUnit)(planUnit),
		}
		todayValue := 0

		for _, score := range habit.HabitScore {
			if isToday(score.DateTime) {
				todayValue += score.Value
			}
		}

		responsHabit := Habit{
			CurrentPlan:   &respPlan,
			Id:            &habit.ID,
			Name:          habit.HabitName,
			Notifications: habit.Notifications,
			StartDate:     &openapi_types.Date{habit.StartDate},
			TodayValue:    todayValue,
		}

		response = append(response, responsHabit)

	}
	return ctx.JSON(http.StatusOK, response)

}

func (h *HabitsApi) PostApiHabits(ctx echo.Context) error {

	var habit NewHabit

	if err := ctx.Bind(&habit); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	userId := ctx.Get("user_id").(int)

	todayValue := new(float32)
	*todayValue = 0

	t := time.Now()

	habit1 := models.Habit{
		HabitName:     habit.Name,
		Description:   habit.Description,
		UserId:        userId,
		StartDate:     t,
		Notifications: habit.Notifications,
	}

	plan := models.Plan{
		HabitId:   habit1.ID,
		PlanUnit:  (*models.PlanUnit)(habit.Plan.PlanUnit),
		Goal:      habit.Plan.Goal,
		StartTime: habit1.StartDate,
	}

	habit1.Plans = append(habit1.Plans, plan)

	if err := h.habitRepo.CreateHabit(&habit1); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	habitResponse := Habit{
		CurrentPlan:   &habit.Plan,
		Id:            &habit1.ID,
		Name:          habit.Name,
		Notifications: habit.Notifications,
		StartDate:     &openapi_types.Date{t},
		TodayValue:    0,
	}

	return ctx.JSON(http.StatusCreated, habitResponse)

}

func (h *HabitsApi) GetApiHabitsHabitId(ctx echo.Context, habitId int) error {
	userId := ctx.Get("user_id").(int)

	habit, err := h.habitRepo.GetHabitById(habitId, userId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var plansArr = habit.Plans
	var plan models.Plan
	if len(plansArr) > 0 {

		sort.Slice(plansArr, func(i, j int) bool {
			return plansArr[i].ID > plansArr[j].ID
		})
		plan = plansArr[0]
	}

	planUnit := plan.PlanUnit

	planResp := Plan{
		Goal:     plan.Goal,
		PlanUnit: (*PlanUnit)(planUnit),
	}

	todayValue := 0

	for _, score := range habit.HabitScore {
		if isToday(score.DateTime) {
			todayValue += score.Value
		}
	}

	response := Habit{
		CurrentPlan:   &planResp,
		Id:            &habit.ID,
		Name:          habit.HabitName,
		Notifications: habit.Notifications,
		StartDate:     &openapi_types.Date{habit.StartDate},
		TodayValue:    todayValue,
	}

	return ctx.JSON(http.StatusOK, response)

}

func (h *HabitsApi) PutApiHabitsHabitId(ctx echo.Context, habitId int) error {

	userId := ctx.Get("user_id").(int)

	var updateHabit HabitUpdate

	if err := ctx.Bind(&updateHabit); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	habit, err := h.habitRepo.GetHabitById(habitId, userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	plans, err := h.planRepo.GetPlansByHabitId(habitId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	var plans_arr = *plans
	var plan models.Plan
	if len(plans_arr) > 0 {

		sort.Slice(plans_arr, func(i, j int) bool {
			return plans_arr[i].ID > plans_arr[j].ID
		})
		plan = plans_arr[0]
	}

	habit.HabitName = *updateHabit.Name
	habit.Notifications = updateHabit.Notifications
	habit.Description = updateHabit.Description

	if err := h.habitRepo.UpdateHabit(habit); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if plan.Goal == updateHabit.Plan.Goal && plan.PlanUnit == (*models.PlanUnit)(updateHabit.Plan.PlanUnit) {
		return ctx.JSON(http.StatusOK, map[string]string{"message": "OK"})
	}

	t := time.Now()

	plan.CloseTime = &t

	if err := h.planRepo.UpdatePlan(&plan); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	plan.Goal = updateHabit.Plan.Goal
	plan.PlanUnit = (*models.PlanUnit)(updateHabit.Plan.PlanUnit)
	plan.HabitId = habitId

	timeStart := time.Now()

	plan.StartTime = timeStart
	plan.CloseTime = nil
	plan.ID = 0

	if err := h.planRepo.CreatePlan(&plan); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "OK"})

}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
