package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_Integration(t *testing.T) {
	// 1. Setup real Postgres container using our helper
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepo(db)
	ctx := context.Background()

	// --- SETUP UNIQUE TEST DATA ---
	// We create these specifically for the tests to avoid ID collisions 
	// and ensure Foreign Key constraints are satisfied.

	cat := model.Category{CategoryName: "Integration Test Category"}
	assert.NoError(t, db.Create(&cat).Error)

	prod := model.Product{
		ProductName: "Integration Test Widget", 
		CategoryID:  cat.ID, 
		UnitPrice:   10.0,
	}
	assert.NoError(t, db.Create(&prod).Error)

	ship := model.Shipper{CompanyName: "Test Logistics"}
	assert.NoError(t, db.Create(&ship).Error)

	emp := model.Employee{EmployeeName: "Testy McTestFace"}
	assert.NoError(t, db.Create(&emp).Error)

	// Use a unique CustomerID (5 chars) that isn't in the CSV (e.g., "ALFKI" is taken)
	cust := model.Customer{CustomerID: "TESTZ", CompanyName: "Test Corp"}
	assert.NoError(t, db.Create(&cust).Error)

	t.Run("Create Order with Details (Deep Insert)", func(t *testing.T) {
		order := &model.Order{
			CustomerID:   cust.CustomerID,
			EmployeeID:   emp.ID, // Must be non-zero
			ShipperID:    ship.ID, // Must be non-zero
			OrderDate:    time.Now(),
			RequiredDate: time.Now().AddDate(0, 0, 7),
			Freight:      15.50,
			OrderDetails: []model.OrderDetail{
				{ProductID: prod.ID, UnitPrice: 10.0, Quantity: 2, Discount: 0},
			},
		}

		err := repo.Create(ctx, order)
		assert.NoError(t, err)
		assert.NotZero(t, order.ID)
		
		if assert.Len(t, order.OrderDetails, 1) {
			assert.NotZero(t, order.OrderDetails[0].ID)
		}
	})

	t.Run("GetByID with Preload", func(t *testing.T) {
		newOrder := &model.Order{
			CustomerID: cust.CustomerID,
			EmployeeID: emp.ID,
			ShipperID:  ship.ID,
			OrderDetails: []model.OrderDetail{
				{ProductID: prod.ID, Quantity: 5, UnitPrice: 10.0},
			},
		}
		assert.NoError(t, db.Create(newOrder).Error)

		found, err := repo.GetByID(ctx, newOrder.ID)
		assert.NoError(t, err)
		
		if assert.NotNil(t, found) && assert.NotEmpty(t, found.OrderDetails) {
			assert.Equal(t, 5, found.OrderDetails[0].Quantity)
			// Verify that the association was preloaded
			assert.Equal(t, cust.CustomerID, found.CustomerID)
		}
	})

	t.Run("GetAll Returns Seeded and New Data", func(t *testing.T) {
		// Northwind CSV contains 830 orders.
		// Our previous tests added at least 2 more.
		orders, err := repo.GetAll(ctx)
		
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(orders), 830)

		// Pick an order and verify Preload ("OrderDetails") worked
		// We check the last one as it's likely one we just created
		lastOrder := orders[len(orders)-1]
		assert.NotEmpty(t, lastOrder.OrderDetails, "OrderDetails must be preloaded in GetAll")
	})

	t.Run("Update with Full Association Save", func(t *testing.T) {
		// Initialize with required Foreign Keys to pass DB constraints
		order := &model.Order{
			CustomerID: cust.CustomerID,
			EmployeeID: emp.ID,
			ShipperID:  ship.ID,
		}
		db.Create(order)

		// Update by adding a detail
		order.OrderDetails = []model.OrderDetail{
			{ProductID: prod.ID, Quantity: 10, UnitPrice: 10.0},
		}
		
		err := repo.Update(ctx, order)
		assert.NoError(t, err)

		// Verify change
		updated, _ := repo.GetByID(ctx, order.ID)
		if assert.NotNil(t, updated) && assert.Len(t, updated.OrderDetails, 1) {
			assert.Equal(t, 10, updated.OrderDetails[0].Quantity)
		}
	})

	t.Run("Delete Cascades Soft-Delete", func(t *testing.T) {
		order := &model.Order{
			CustomerID: cust.CustomerID,
			EmployeeID: emp.ID,
			ShipperID:  ship.ID,
			OrderDetails: []model.OrderDetail{
				{ProductID: prod.ID, Quantity: 1, UnitPrice: 10.0},
			},
		}
		db.Create(order)
		
		if len(order.OrderDetails) == 0 {
			t.Fatal("Failed to setup order details for delete test")
		}
		detailID := order.OrderDetails[0].ID

		err := repo.Delete(ctx, order.ID)
		assert.NoError(t, err)

		// 1. Verify Order is soft-deleted (standard query fails)
		var checkOrder model.Order
		err = db.First(&checkOrder, order.ID).Error
		assert.Error(t, err)

		// 2. Verify OrderDetail is also soft-deleted (Cascading)
		var checkDetail model.OrderDetail
		err = db.First(&checkDetail, detailID).Error
		assert.Error(t, err)
		
		// 3. Physically check the DB using Unscoped
		db.Unscoped().First(&checkDetail, detailID)
		assert.NotNil(t, checkDetail.DeletedAt, "DeletedAt should be populated")
	})
}