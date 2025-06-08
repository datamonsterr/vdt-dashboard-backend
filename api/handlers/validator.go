package handlers

import (
	"net/http"

	"vdt-dashboard-backend/models"
	"vdt-dashboard-backend/services"

	"github.com/gin-gonic/gin"
)

// ValidatorHandler handles validation requests
type ValidatorHandler struct {
	validatorService    services.ValidatorService
	sqlGeneratorService services.SQLGeneratorService
}

// NewValidatorHandler creates a new validator handler
func NewValidatorHandler(validatorService services.ValidatorService, sqlGeneratorService services.SQLGeneratorService) *ValidatorHandler {
	return &ValidatorHandler{
		validatorService:    validatorService,
		sqlGeneratorService: sqlGeneratorService,
	}
}

// ValidateSchema handles POST /schemas/validate
func (h *ValidatorHandler) ValidateSchema(c *gin.Context) {
	var request models.SchemaValidationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request data", models.ErrValidation, err.Error()))
		return
	}

	validationResult, err := h.validatorService.ValidateSchema(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Validation failed", models.ErrInternalError, err.Error()))
		return
	}

	// If validation passed, generate SQL preview
	if validationResult.Valid {
		schemaData := models.SchemaData{
			Tables:      request.Tables,
			ForeignKeys: request.ForeignKeys,
		}

		sqlStatements, err := h.sqlGeneratorService.GenerateCreateTables(schemaData)
		if err == nil {
			validationResult.GeneratedSQL = sqlStatements
		}
	}

	statusCode := http.StatusOK
	message := "Schema is valid"

	if !validationResult.Valid {
		statusCode = http.StatusBadRequest
		message = "Schema validation failed"
	}

	c.JSON(statusCode, models.SuccessResponse(message, validationResult))
}
