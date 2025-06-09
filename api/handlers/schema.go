package handlers

import (
	"net/http"

	"vdt-dashboard-backend/api/middleware"
	"vdt-dashboard-backend/models"
	"vdt-dashboard-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SchemaHandler handles schema-related HTTP requests
type SchemaHandler struct {
	schemaService services.SchemaService
}

// NewSchemaHandler creates a new schema handler
func NewSchemaHandler(schemaService services.SchemaService) *SchemaHandler {
	return &SchemaHandler{
		schemaService: schemaService,
	}
}

// CreateSchema handles POST /schemas
func (h *SchemaHandler) CreateSchema(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	var request models.CreateSchemaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request data", models.ErrValidation, err.Error()))
		return
	}

	schema, err := h.schemaService.CreateSchema(request, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create schema", models.ErrInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Schema created successfully", schema))
}

// ListSchemas handles GET /schemas
func (h *SchemaHandler) ListSchemas(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	var pagination models.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid pagination parameters", models.ErrValidation, err.Error()))
		return
	}

	schemas, paginationResp, err := h.schemaService.ListSchemas(pagination, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to list schemas", models.ErrInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.PaginatedSuccessResponse("Schemas retrieved successfully", schemas, paginationResp))
}

// GetSchema handles GET /schemas/:id
func (h *SchemaHandler) GetSchema(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid schema ID", models.ErrValidation, "ID must be a valid UUID"))
		return
	}

	schema, err := h.schemaService.GetSchema(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Schema not found", models.ErrSchemaNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Schema retrieved successfully", schema))
}

// UpdateSchema handles PUT /schemas/:id
func (h *SchemaHandler) UpdateSchema(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid schema ID", models.ErrValidation, "ID must be a valid UUID"))
		return
	}

	var request models.UpdateSchemaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request data", models.ErrValidation, err.Error()))
		return
	}

	schema, err := h.schemaService.UpdateSchema(id, userID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to update schema", models.ErrInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Schema updated successfully", schema))
}

// DeleteSchema handles DELETE /schemas/:id
func (h *SchemaHandler) DeleteSchema(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid schema ID", models.ErrValidation, "ID must be a valid UUID"))
		return
	}

	if err := h.schemaService.DeleteSchema(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to delete schema", models.ErrInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Schema deleted successfully", gin.H{"id": id}))
}

// ExportSQL handles GET /schemas/:id/export/sql
func (h *SchemaHandler) ExportSQL(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated", models.ErrUnauthorized, "Missing user context"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid schema ID", models.ErrValidation, "ID must be a valid UUID"))
		return
	}

	sqlExport, err := h.schemaService.ExportSQL(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to export SQL", models.ErrInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("SQL export generated", sqlExport))
}
