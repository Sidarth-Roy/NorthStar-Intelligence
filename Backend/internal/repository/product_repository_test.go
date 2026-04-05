package repository

import (
	"context"
	"testing"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository_Integration(t *testing.T) {
	// 1. Setup the real Postgres container
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewProductRepo(db)
	ctx := context.Background()

	// 2. Create a Category first (Product needs a CategoryID)
	testCategory := model.Category{
		CategoryName: "Test Category",
		Description:  "Description for testing",
	}
	db.Create(&testCategory)

	t.Run("Create and Get Product", func(t *testing.T) {
		product := &model.Product{
			ProductName:     "Test Product",
			QuantityPerUnit: "10 boxes",
			UnitPrice:       19.99,
			CategoryID:      testCategory.ID, // Using the ID we just created
		}

		// Test Create
		err := repo.Create(ctx, product)
		assert.NoError(t, err)
		assert.NotZero(t, product.ID)

		// Test GetByID
		found, err := repo.GetByID(ctx, product.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Test Product", found.ProductName)
		assert.Equal(t, 19.99, found.UnitPrice)
	})

	t.Run("GetAll Products", func(t *testing.T) {
		// Add another product
		repo.Create(ctx, &model.Product{
			ProductName: "Another One",
			CategoryID:  testCategory.ID,
		})

		products, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		// Expecting at least 2 products based on the previous sub-test
		assert.GreaterOrEqual(t, len(products), 2)
	})

	t.Run("Update Product Details", func(t *testing.T) {
		product := &model.Product{ProductName: "Old Name", CategoryID: testCategory.ID}
		repo.Create(ctx, product)

		product.ProductName = "New Shiny Name"
		product.UnitPrice = 99.99
		
		err := repo.Update(ctx, product)
		assert.NoError(t, err)

		found, _ := repo.GetByID(ctx, product.ID)
		assert.Equal(t, "New Shiny Name", found.ProductName)
		assert.Equal(t, 99.99, found.UnitPrice)
	})

	t.Run("Soft Delete Product", func(t *testing.T) {
		product := &model.Product{ProductName: "Delete Me", CategoryID: testCategory.ID}
		repo.Create(ctx, product)

		err := repo.Delete(ctx, product.ID)
		assert.NoError(t, err)

		// Verify it's "gone" from normal queries
		_, err = repo.GetByID(ctx, product.ID)
		assert.Error(t, err) // Should be record not found

		// Verify it's still in DB (Soft Delete check)
		var deleted model.Product
		db.Unscoped().First(&deleted, product.ID)
		assert.NotNil(t, deleted.DeletedAt)
	})
}