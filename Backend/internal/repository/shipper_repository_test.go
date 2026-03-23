package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestShipperRepository_Integration(t *testing.T) {
	// 1. Setup real database container using our global helper
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewShipperRepo(db)
	ctx := context.Background()

	t.Run("Create and Find Shipper", func(t *testing.T) {
		shipper := &model.Shipper{
			CompanyName: "Speedy Express",
		}

		// Test Create
		err := repo.Create(ctx, shipper)
		assert.NoError(t, err)
		assert.NotZero(t, shipper.ID)

		// Test GetByID
		found, err := repo.GetByID(ctx, shipper.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Speedy Express", found.CompanyName)
	})

	t.Run("GetAll Returns List", func(t *testing.T) {
		// Add another shipper to the database
		repo.Create(ctx, &model.Shipper{CompanyName: "United Package"})
		repo.Create(ctx, &model.Shipper{CompanyName: "Federal Shipping"})

		shippers, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		// We expect at least 3 (1 from previous test + 2 new ones)
		assert.GreaterOrEqual(t, len(shippers), 3)
	})

	t.Run("Update Shipper Data", func(t *testing.T) {
		shipper := &model.Shipper{CompanyName: "Old Name"}
		repo.Create(ctx, shipper)

		// Modify
		shipper.CompanyName = "New Global Logistics"
		err := repo.Update(ctx, shipper)
		assert.NoError(t, err)

		// Verify change in DB
		found, _ := repo.GetByID(ctx, shipper.ID)
		assert.Equal(t, "New Global Logistics", found.CompanyName)
	})

	t.Run("Verify Soft Delete", func(t *testing.T) {
		shipper := &model.Shipper{CompanyName: "Temporary Shipper"}
		repo.Create(ctx, shipper)
		id := shipper.ID

		// Perform Delete
		err := repo.Delete(ctx, id)
		assert.NoError(t, err)

		// 1. Standard GetByID should fail (GORM filters out soft-deleted)
		_, err = repo.GetByID(ctx, id)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// 2. Verification: Check if record exists physically in DB (Unscoped)
		var deleted model.Shipper
		db.Unscoped().First(&deleted, id)
		assert.NotNil(t, deleted.DeletedAt, "DeletedAt field should not be nil")
		assert.Equal(t, "Temporary Shipper", deleted.CompanyName)
	})
}	