package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCustomerRepository_Integration(t *testing.T) {
	// 1. Setup real Postgres container using our helper
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewCustomerRepo(db)
	ctx := context.Background()

	t.Run("Create and Get Customer", func(t *testing.T) {
		customer := &model.Customer{
			CustomerID:   "NEW12",
			CompanyName:  "Integration Test Corp",
			ContactName:  "Test Contact",
			ContactTitle: "Test Manager",
			City:         "Test City",
			Country:      "Testland",
		}

		// Test Create
		err := repo.Create(ctx, customer)
		assert.NoError(t, err)
		assert.NotZero(t, customer.ID) // Auto-increment ID from Base model

		// Test GetByID
		found, err := repo.GetByID(ctx, customer.ID)
		assert.NoError(t, err)
		assert.Equal(t, "NEW12", found.CustomerID)
		assert.Equal(t, "Test City", found.City)
	})

	t.Run("Get All Customers", func(t *testing.T) {
		// Insert another one
		repo.Create(ctx, &model.Customer{
			CustomerID:  "ANATR",
			CompanyName: "Ana Trujillo Emparedados",
		})

		customers, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		// Expect at least 2 because of the previous sub-test
		assert.GreaterOrEqual(t, len(customers), 2)
	})

	t.Run("Update Customer Details", func(t *testing.T) {
		customer := &model.Customer{CustomerID: "EDIT1", CompanyName: "Original Corp"}
		repo.Create(ctx, customer)

		customer.CompanyName = "Updated Corp"
		customer.City = "London"
		
		err := repo.Update(ctx, customer)
		assert.NoError(t, err)

		found, _ := repo.GetByID(ctx, customer.ID)
		assert.Equal(t, "Updated Corp", found.CompanyName)
		assert.Equal(t, "London", found.City)
	})

	t.Run("Soft Delete Customer", func(t *testing.T) {
		customer := &model.Customer{CustomerID: "DEL01", CompanyName: "Gone Soon"}
		repo.Create(ctx, customer)

		err := repo.Delete(ctx, customer.ID)
		assert.NoError(t, err)

		// Verification: Normal query should not find it
		_, err = repo.GetByID(ctx, customer.ID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// Verification: Check DB directly to ensure row still exists but DeletedAt is set
		var deleted model.Customer
		db.Unscoped().First(&deleted, customer.ID)
		assert.NotNil(t, deleted.DeletedAt)
	})
}