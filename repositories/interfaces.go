package repositories

import (
	"vdt-dashboard-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SchemaRepository defines the interface for schema data access
type SchemaRepository interface {
	Create(schema *models.Schema) error
	GetByID(id uuid.UUID) (*models.Schema, error)
	GetByIDAndUserID(id, userID uuid.UUID) (*models.Schema, error)
	GetByName(name string) (*models.Schema, error)
	GetByNameAndUserID(name string, userID uuid.UUID) (*models.Schema, error)
	List(pagination models.PaginationRequest) ([]models.SchemaListResponse, int, error)
	ListByUserID(pagination models.PaginationRequest, userID uuid.UUID) ([]models.SchemaListResponse, int, error)
	Update(schema *models.Schema) error
	Delete(id uuid.UUID) error
	DeleteByIDAndUserID(id, userID uuid.UUID) error
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByClerkID(clerkID string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

// NewSchemaRepository creates a new schema repository
func NewSchemaRepository(db *gorm.DB) SchemaRepository {
	return &schemaRepository{db: db}
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// schemaRepository implements SchemaRepository
type schemaRepository struct {
	db *gorm.DB
}

// Create creates a new schema
func (r *schemaRepository) Create(schema *models.Schema) error {
	return r.db.Create(schema).Error
}

// GetByID gets a schema by ID
func (r *schemaRepository) GetByID(id uuid.UUID) (*models.Schema, error) {
	var schema models.Schema
	err := r.db.Where("id = ?", id).First(&schema).Error
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// GetByIDAndUserID gets a schema by ID and user ID
func (r *schemaRepository) GetByIDAndUserID(id, userID uuid.UUID) (*models.Schema, error) {
	var schema models.Schema
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&schema).Error
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// GetByName gets a schema by name
func (r *schemaRepository) GetByName(name string) (*models.Schema, error) {
	var schema models.Schema
	err := r.db.Where("name = ?", name).First(&schema).Error
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// GetByNameAndUserID gets a schema by name and user ID
func (r *schemaRepository) GetByNameAndUserID(name string, userID uuid.UUID) (*models.Schema, error) {
	var schema models.Schema
	err := r.db.Where("name = ? AND user_id = ?", name, userID).First(&schema).Error
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// List gets paginated list of schemas
func (r *schemaRepository) List(pagination models.PaginationRequest) ([]models.SchemaListResponse, int, error) {
	var schemas []models.Schema
	var total int64

	query := r.db.Model(&models.Schema{})

	// Add search filter if provided
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	if err := query.Offset(offset).Limit(pagination.Limit).Find(&schemas).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response format
	var response []models.SchemaListResponse
	for _, schema := range schemas {
		// Safely get table count - handle case where SchemaDefinition.Tables might be nil
		tableCount := 0
		if schema.SchemaDefinition.Tables != nil {
			tableCount = len(schema.SchemaDefinition.Tables)
		}

		response = append(response, models.SchemaListResponse{
			ID:           schema.ID,
			Name:         schema.Name,
			Description:  schema.Description,
			DatabaseName: schema.DatabaseName,
			Status:       schema.Status,
			TableCount:   tableCount,
			CreatedAt:    schema.CreatedAt,
			UpdatedAt:    schema.UpdatedAt,
			Version:      schema.Version,
		})
	}

	return response, int(total), nil
}

// ListByUserID gets paginated list of schemas for a specific user
func (r *schemaRepository) ListByUserID(pagination models.PaginationRequest, userID uuid.UUID) ([]models.SchemaListResponse, int, error) {
	var schemas []models.Schema
	var total int64

	query := r.db.Model(&models.Schema{}).Where("user_id = ?", userID)

	// Add search filter if provided
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	if err := query.Offset(offset).Limit(pagination.Limit).Find(&schemas).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response format
	var response []models.SchemaListResponse
	for _, schema := range schemas {
		// Safely get table count - handle case where SchemaDefinition.Tables might be nil
		tableCount := 0
		if schema.SchemaDefinition.Tables != nil {
			tableCount = len(schema.SchemaDefinition.Tables)
		}

		response = append(response, models.SchemaListResponse{
			ID:           schema.ID,
			Name:         schema.Name,
			Description:  schema.Description,
			DatabaseName: schema.DatabaseName,
			Status:       schema.Status,
			TableCount:   tableCount,
			CreatedAt:    schema.CreatedAt,
			UpdatedAt:    schema.UpdatedAt,
			Version:      schema.Version,
		})
	}

	return response, int(total), nil
}

// Update updates a schema
func (r *schemaRepository) Update(schema *models.Schema) error {
	return r.db.Save(schema).Error
}

// Delete soft deletes a schema
func (r *schemaRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Schema{}).Error
}

// DeleteByIDAndUserID soft deletes a schema by ID and user ID
func (r *schemaRepository) DeleteByIDAndUserID(id, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Schema{}).Error
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByClerkID gets a user by Clerk ID
func (r *userRepository) GetByClerkID(clerkID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("clerk_user_id = ?", clerkID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}
