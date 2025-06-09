package handlers

import (
	"net/http"

	"vdt-dashboard-backend/api/middleware"
	"vdt-dashboard-backend/models"
	"vdt-dashboard-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DatabaseHandler handles database management requests
type DatabaseHandler struct {
	databaseManagerService services.DatabaseManagerService
	schemaService          services.SchemaService
}

// NewDatabaseHandler creates a new database handler
func NewDatabaseHandler(databaseManagerService services.DatabaseManagerService, schemaService services.SchemaService) *DatabaseHandler {
	return &DatabaseHandler{
		databaseManagerService: databaseManagerService,
		schemaService:          schemaService,
	}
}

// GetDatabaseStatus handles GET /schemas/:id/database/status
func (h *DatabaseHandler) GetDatabaseStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid schema ID", models.ErrValidation, "ID must be a valid UUID"))
		return
	}
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	schema, err := h.schemaService.GetSchema(id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Schema not found", models.ErrSchemaNotFound, err.Error()))
		return
	}

	status, err := h.databaseManagerService.GetDatabaseStatus(schema.DatabaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to get database status", models.ErrDatabaseError, err.Error()))
		return
	}

	status.SchemaID = schema.ID

	c.JSON(http.StatusOK, models.SuccessResponse("Database status retrieved", status))
}

// RegenerateDatabase handles POST /schemas/:id/database/regenerate
func (h *DatabaseHandler) RegenerateDatabase(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid schema ID", models.ErrValidation, "ID must be a valid UUID"))
		return
	}

	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	schema, err := h.schemaService.GetSchema(id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Schema not found", models.ErrSchemaNotFound, err.Error()))
		return
	}

	err = h.databaseManagerService.RegenerateDatabase(schema.SchemaDefinition, schema.DatabaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to regenerate database", models.ErrDatabaseError, err.Error()))
		return
	}

	response := gin.H{
		"schemaId":      schema.ID,
		"databaseName":  schema.DatabaseName,
		"status":        "regenerated",
		"regeneratedAt": "2024-01-01T12:30:00Z", // TODO: Use actual timestamp
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Database regenerated successfully", response))
}
