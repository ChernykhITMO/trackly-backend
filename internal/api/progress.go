package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"trackly-backend/internal/models"
	"trackly-backend/internal/repositories"
)

type ProgressApi struct {
	statisticRepo *repositories.StatisticRepository
	habitRepo     *repositories.HabitRepository
	planRepo      *repositories.PlanRepository
}

func NewProgressApi(statisticRepo *repositories.StatisticRepository, habitRepo *repositories.HabitRepository, planRepo *repositories.PlanRepository) *ProgressApi {
	return &ProgressApi{statisticRepo: statisticRepo, habitRepo: habitRepo, planRepo: planRepo}
}

func (h *ProgressApi) PostApiHabitsHabitIdScore(ctx echo.Context, habitId int) error {

	var updatedScore ScoreHabit

	userId := ctx.Get("user_id").(int)

	if err := ctx.Bind(&updatedScore); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if habit, _ := h.habitRepo.GetHabitById(habitId, userId); habit == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Habit not found"})
	}

	plans, err := h.planRepo.GetPlansByHabitId(habitId)
	if err != nil {
		return ctx.JSON(500, map[string]string{"message": err.Error()})
	}

	var plans_arr = *plans
	var plan models.Plan
	if len(plans_arr) > 0 {

		sort.Slice(plans_arr, func(i, j int) bool {
			return plans_arr[i].ID > plans_arr[j].ID
		})
		plan = plans_arr[0]
	}

	score := models.HabitScore{
		HabitId:  plan.HabitId,
		PlanId:   plan.ID,
		DateTime: updatedScore.Date,
		Value:    updatedScore.Value,
	}

	if err := h.statisticRepo.CreateStatistic(&score); err != nil {
		return ctx.JSON(500, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(200, map[string]string{"message": "Habit score created"})

}
