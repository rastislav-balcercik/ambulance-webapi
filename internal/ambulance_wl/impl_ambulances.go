package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rastislav-balcercik/ambulance-webapi/internal/db_service"
)

// Kópia zakomentovanej časti z api_ambulances.go
// CreateAmbulance - Saves new ambulance definition
func (this *implAmbulancesAPI) CreateAmbulance(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	ambulance := Ambulance{}
	err := ctx.BindJSON(&ambulance)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if ambulance.Id == "" {
		ambulance.Id = uuid.New().String()
	}

	err = db.CreateDocument(ctx, ambulance.Id, &ambulance)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusCreated,
			ambulance,
		)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Ambulance already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create ambulance in database",
				"error":   err.Error(),
			},
		)
	}
}

// DeleteAmbulance - Deletes specific ambulance
func (this *implAmbulancesAPI) DeleteAmbulance(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	ambulanceId := ctx.Param("ambulanceId")
	err := db.DeleteDocument(ctx, ambulanceId)

	switch err {
	case nil:
		ctx.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ambulance not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete ambulance from database",
				"error":   err.Error(),
			})
	}
}

// GetAmbulance - Retrieves specific ambulance
func (this *implAmbulancesAPI) GetAmbulance(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	ambulanceId := ctx.Param("ambulanceId")
	ambulance, err := db.FindDocument(ctx, ambulanceId)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusOK,
			ambulance,
		)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ambulance not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete ambulance from database",
				"error":   err.Error(),
			})
	}
}

// UpdateAmbulance - Updates specific amulance entry
func (this *implAmbulancesAPI) UpdateAmbulance(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	// Extract ambulance ID from URL parameter
	ambulanceId := ctx.Param("ambulanceId")

	// Check if the ambulance exists
	_, err := db.FindDocument(ctx, ambulanceId)
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ambulance not found",
				"error":   err.Error(),
			},
		)
		return
	}

	// Bind JSON data to Ambulance struct
	var updatedAmbulance Ambulance
	if err := ctx.BindJSON(&updatedAmbulance); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	// Update ambulance in the database
	err = db.UpdateDocument(ctx, ambulanceId, &updatedAmbulance)
	if err != nil {
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update ambulance in database",
				"error":   err.Error(),
			})
		return
	}

	// Return updated ambulance details
	ctx.JSON(
		http.StatusOK,
		updatedAmbulance,
	)
}
