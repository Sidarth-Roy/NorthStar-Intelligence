package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCategoryRepository_Integration(t *testing.T) {
	// Setup real database container
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewCategoryRepo(db)
	ctx := context.Background()

	t.Run("Create and Find Category", func(t *testing.T) {
		cat := &model.Category{
			CategoryName: "Testing-Category-999", // Change from "Hardware"
			Description:  "Integration Test",
		}

		err := repo.Create(ctx, cat)
		assert.NoError(t, err)
		assert.NotZero(t, cat.ID) // This will now be 9 (since CSV has 8 categories)

		found, err := repo.GetByID(ctx, cat.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Testing-Category-999", found.CategoryName)
	})

	t.Run("GetAll Returns All Records", func(t *testing.T) {
		// Insert another record
		// repo.Create(ctx, &model.Category{CategoryName: "Software"})

		categories, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(categories), 2)
	})

	t.Run("Update Existing Category", func(t *testing.T) {
		cat := &model.Category{CategoryName: "Original"}
		repo.Create(ctx, cat)

		cat.CategoryName = "Updated Name"
		err := repo.Update(ctx, cat)
		assert.NoError(t, err)

		found, _ := repo.GetByID(ctx, cat.ID)
		assert.Equal(t, "Updated Name", found.CategoryName)
	})

	t.Run("Soft Delete Logic", func(t *testing.T) {
		cat := &model.Category{CategoryName: "To Be Deleted"}
		repo.Create(ctx, cat)

		err := repo.Delete(ctx, cat.ID)
		assert.NoError(t, err)

		// GetByID should fail (standard query ignores soft-deleted)
		_, err = repo.GetByID(ctx, cat.ID)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// Check if it exists in DB but with DeletedAt set
		var deleted model.Category
		db.Unscoped().First(&deleted, cat.ID)
		assert.NotNil(t, deleted.DeletedAt)
	})
}