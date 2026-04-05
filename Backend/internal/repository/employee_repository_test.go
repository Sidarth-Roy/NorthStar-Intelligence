package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestEmployeeRepository_Integration(t *testing.T) {
	// 1. Setup real Postgres container
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewEmployeeRepo(db)
	ctx := context.Background()

	t.Run("Create and Find Employee", func(t *testing.T) {
		emp := &model.Employee{
			EmployeeName: "Test Robot",
			Title:        "Quality Assurance",
			City:         "Seattle",
			Country:      "USA",
		}

		err := repo.Create(ctx, emp)
		assert.NoError(t, err)
		assert.NotZero(t, emp.ID)

		found, err := repo.GetByID(ctx, emp.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, found) { // Safety check to prevent panic
        	assert.Equal(t, "Test Robot", found.EmployeeName)
    	}
	})

	t.Run("Self-Referencing Relationship (ReportsTo)", func(t *testing.T) {
		// Create Manager
		manager := &model.Employee{EmployeeName: "Andrew Fuller", Title: "Vice President"}
		repo.Create(ctx, manager)

		// Create Subordinate reporting to Manager
		subordinate := &model.Employee{
			EmployeeName: "Anne Dodsworth",
			ReportsTo:    &manager.ID, // Link to Manager
		}
		err := repo.Create(ctx, subordinate)
		assert.NoError(t, err)

		// Verify relationship
		found, _ := repo.GetByID(ctx, subordinate.ID)
		assert.NotNil(t, found.ReportsTo)
		assert.Equal(t, manager.ID, *found.ReportsTo)
	})

	t.Run("GetAll Employees", func(t *testing.T) {
		// We already added 3 employees in previous sub-tests within this container
		employees, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(employees), 2)
	})

	t.Run("Update Employee Details", func(t *testing.T) {
		emp := &model.Employee{EmployeeName: "Janet Leverling", City: "London"}
		repo.Create(ctx, emp)

		emp.City = "Kirkland"
		err := repo.Update(ctx, emp)
		assert.NoError(t, err)

		found, _ := repo.GetByID(ctx, emp.ID)
		assert.Equal(t, "Kirkland", found.City)
	})

	t.Run("Delete (Soft Delete) Verification", func(t *testing.T) {
		emp := &model.Employee{EmployeeName: "Steven Buchanan"}
		repo.Create(ctx, emp)

		err := repo.Delete(ctx, emp.ID)
		assert.NoError(t, err)

		// Standard query should not find it
		_, err = repo.GetByID(ctx, emp.ID)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// Unscoped query should find the record with DeletedAt set
		var deleted model.Employee
		db.Unscoped().First(&deleted, emp.ID)
		assert.NotNil(t, deleted.DeletedAt)
	})
}