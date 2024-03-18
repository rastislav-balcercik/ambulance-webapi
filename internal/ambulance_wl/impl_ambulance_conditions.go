package ambulance_wl

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_conditions.go
func (this *implAmbulanceConditionsAPI) GetConditions(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(
		ctx *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		result := ambulance.PredefinedConditions
		if result == nil {
			result = []Condition{}
		}
		return nil, result, http.StatusOK
	})
}

func (this *implAmbulanceConditionsAPI) CreateCondition(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(
		ctx *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		var condition Condition

		if err := ctx.ShouldBindJSON(&condition); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if ambulance.Id == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Ambulance ID is required",
			}, http.StatusBadRequest
		}

		conflictIndx := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return c.Code == condition.Code
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		ambulance.PredefinedConditions = append(ambulance.PredefinedConditions, condition)
		// entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return c.Code == condition.Code
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}
		return ambulance, ambulance.PredefinedConditions[entryIndx], http.StatusOK
	})
}

func (this *implAmbulanceConditionsAPI) UpdateCondition(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var condition Condition

		if err := c.ShouldBindJSON(&condition); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		conditionCode := ctx.Param("conditionId")

		if conditionCode == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Condition code is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return conditionCode == c.Code
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		entryWithSameId := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return condition.Code == c.Code
		})

		if entryWithSameId > 0 {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "You cannot change the code of the condition, because conditon with that code already exists",
			}, http.StatusBadRequest
		}

		if condition.Code != "" {
			ambulance.PredefinedConditions[entryIndx].Code = condition.Code
		}

		if condition.Reference != "" {
			ambulance.PredefinedConditions[entryIndx].Reference = condition.Reference
		}

		ambulance.PredefinedConditions[entryIndx].TypicalDurationMinutes = condition.TypicalDurationMinutes

		if condition.Value != "" {
			ambulance.PredefinedConditions[entryIndx].Value = condition.Value
		}

		return ambulance, ambulance.PredefinedConditions[entryIndx], http.StatusOK
	})
}

func (this *implAmbulanceConditionsAPI) DeleteCondition(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		conditionCode := ctx.Param("conditionId")

		if conditionCode == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Condition code is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return conditionCode == c.Code
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Entry not found",
			}, http.StatusNotFound
		}

		ambulance.PredefinedConditions = append(ambulance.PredefinedConditions[:entryIndx], ambulance.PredefinedConditions[entryIndx+1:]...)

		return ambulance, nil, http.StatusNoContent
	})
}
