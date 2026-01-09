package repositories

import (
	"testing"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewOrganizationRepository(db)

	t.Run("create organization successfully", func(t *testing.T) {
		req := &models.CreateOrganizationRequest{
			Name: "Test Org",
			Slug: "testorg",
		}

		org, err := repo.Create(req)
		require.NoError(t, err)
		assert.NotZero(t, org.ID)
		assert.Equal(t, "Test Org", org.Name)
		assert.Equal(t, "testorg", org.Slug)
	})

	t.Run("fail on duplicate slug", func(t *testing.T) {
		req1 := &models.CreateOrganizationRequest{Name: "Org 1", Slug: "duplicateslug"}
		req2 := &models.CreateOrganizationRequest{Name: "Org 2", Slug: "duplicateslug"}

		_, err := repo.Create(req1)
		require.NoError(t, err)

		_, err = repo.Create(req2)
		assert.Error(t, err)
	})
}

func TestOrganizationRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewOrganizationRepository(db)

	t.Run("get existing organization", func(t *testing.T) {
		org := testutils.CreateTestOrganization(t, db, "Find Me", "find-me")

		found, err := repo.GetByID(org.ID)
		require.NoError(t, err)
		assert.Equal(t, org.ID, found.ID)
		assert.Equal(t, "Find Me", found.Name)
	})

	t.Run("return error for non-existent organization", func(t *testing.T) {
		found, err := repo.GetByID(99999)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestOrganizationRepository_GetBySlug(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewOrganizationRepository(db)

	t.Run("get organization by slug", func(t *testing.T) {
		org := testutils.CreateTestOrganization(t, db, "Slug Test", "slug-test")

		found, err := repo.GetBySlug("slug-test")
		require.NoError(t, err)
		assert.Equal(t, org.ID, found.ID)
	})

	t.Run("return error for non-existent slug", func(t *testing.T) {
		found, err := repo.GetBySlug("does-not-exist")
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestOrganizationRepository_List(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewOrganizationRepository(db)

	// Create test organizations
	for i := 1; i <= 5; i++ {
		testutils.CreateTestOrganization(t, db, "Org "+string(rune('0'+i)), "org-"+string(rune('0'+i)))
	}

	t.Run("list with limit and offset", func(t *testing.T) {
		orgs, err := repo.List(3, 0)
		require.NoError(t, err)
		assert.Len(t, orgs, 3)
	})

	t.Run("list with offset", func(t *testing.T) {
		orgs, err := repo.List(10, 2)
		require.NoError(t, err)
		assert.Len(t, orgs, 3) // 5 total - 2 offset = 3
	})
}

func TestOrganizationRepository_Update(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewOrganizationRepository(db)

	t.Run("update organization", func(t *testing.T) {
		org := testutils.CreateTestOrganization(t, db, "Original", "original")

		newName := "Updated Name"
		isActive := false
		updates := &models.UpdateOrganizationRequest{
			Name:     &newName,
			IsActive: &isActive,
		}

		err := repo.Update(org.ID, updates)
		require.NoError(t, err)

		found, _ := repo.GetByID(org.ID)
		assert.Equal(t, "Updated Name", found.Name)
		assert.False(t, found.IsActive)
	})
}

func TestOrganizationRepository_Delete(t *testing.T) {
	db, cleanup := testutils.SetupTestDBFile(t)
	defer cleanup()

	repo := NewOrganizationRepository(db)

	t.Run("soft delete organization", func(t *testing.T) {
		org := testutils.CreateTestOrganization(t, db, "Delete Me", "deleteme")

		err := repo.Delete(org.ID)
		require.NoError(t, err)

		// Soft delete should set IsActive to false
		found, _ := repo.GetByID(org.ID)
		assert.NotNil(t, found)
		assert.False(t, found.IsActive)
	})
}
