package services

import (
	"errors"
	"testing"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationService_Create(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	req := &models.CreateOrganizationRequest{
		Name: "Test Org",
		Slug: "test-org",
	}

	org, err := service.Create(req)
	require.NoError(t, err)
	assert.Equal(t, "Test Org", org.Name)
	assert.Equal(t, "test-org", org.Slug)
	assert.Equal(t, int64(1), org.ID)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "organization.created", emitter.EmittedEvents[0].EventType)
}

func TestOrganizationService_Create_RepoError(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	repo.CreateError = errors.New("database error")
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	_, err := service.Create(&models.CreateOrganizationRequest{
		Name: "Test Org",
		Slug: "test-org",
	})
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestOrganizationService_GetByID(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	// Create an org first
	created, _ := repo.Create(&models.CreateOrganizationRequest{
		Name: "Test Org",
		Slug: "test-org",
	})

	org, err := service.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Test Org", org.Name)
}

func TestOrganizationService_GetByID_NotFound(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	org, err := service.GetByID(999)
	require.NoError(t, err)
	assert.Nil(t, org)
}

func TestOrganizationService_GetBySlug(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	// Create an org first
	repo.Create(&models.CreateOrganizationRequest{
		Name: "Test Org",
		Slug: "test-org",
	})

	org, err := service.GetBySlug("test-org")
	require.NoError(t, err)
	assert.Equal(t, "Test Org", org.Name)
	assert.Equal(t, "test-org", org.Slug)
}

func TestOrganizationService_GetBySlug_NotFound(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	org, err := service.GetBySlug("non-existent")
	require.NoError(t, err)
	assert.Nil(t, org)
}

func TestOrganizationService_List(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	// Create multiple orgs
	repo.Create(&models.CreateOrganizationRequest{Name: "Org 1", Slug: "org-1"})
	repo.Create(&models.CreateOrganizationRequest{Name: "Org 2", Slug: "org-2"})

	orgs, err := service.List(10, 0)
	require.NoError(t, err)
	assert.Len(t, orgs, 2)
}

func TestOrganizationService_List_DefaultLimit(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	// Create one org
	repo.Create(&models.CreateOrganizationRequest{Name: "Org 1", Slug: "org-1"})

	// Test with invalid limit (should default to 20)
	orgs, err := service.List(0, 0)
	require.NoError(t, err)
	assert.NotNil(t, orgs)

	// Test with limit > 100 (should default to 20)
	orgs, err = service.List(200, 0)
	require.NoError(t, err)
	assert.NotNil(t, orgs)
}

func TestOrganizationService_Update(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	// Create an org first
	created, _ := repo.Create(&models.CreateOrganizationRequest{
		Name: "Test Org",
		Slug: "test-org",
	})

	newName := "Updated Org"
	err := service.Update(created.ID, &models.UpdateOrganizationRequest{
		Name: &newName,
	})
	require.NoError(t, err)

	// Verify the update
	org, _ := repo.GetByID(created.ID)
	assert.Equal(t, "Updated Org", org.Name)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "organization.updated", emitter.EmittedEvents[0].EventType)
}

func TestOrganizationService_Update_RepoError(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	repo.UpdateError = errors.New("update failed")
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	newName := "Updated Org"
	err := service.Update(1, &models.UpdateOrganizationRequest{
		Name: &newName,
	})
	assert.Error(t, err)
}

func TestOrganizationService_Delete(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	// Create an org first
	created, _ := repo.Create(&models.CreateOrganizationRequest{
		Name: "Test Org",
		Slug: "test-org",
	})

	err := service.Delete(created.ID)
	require.NoError(t, err)

	// Verify deletion
	org, _ := repo.GetByID(created.ID)
	assert.Nil(t, org)

	// Wait for async event emission
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, emitter.EmittedEvents, 1)
	assert.Equal(t, "organization.deleted", emitter.EmittedEvents[0].EventType)
}

func TestOrganizationService_Delete_RepoError(t *testing.T) {
	repo := testutils.NewMockOrganizationRepository()
	repo.DeleteError = errors.New("delete failed")
	emitter := testutils.NewMockEmitter()
	service := NewOrganizationService(repo, emitter)

	err := service.Delete(1)
	assert.Error(t, err)
	assert.Equal(t, "delete failed", err.Error())
}
