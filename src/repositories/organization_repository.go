package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
)

type OrganizationRepository interface {
	Create(org *models.CreateOrganizationRequest) (*models.Organization, error)
	GetByID(id int64) (*models.Organization, error)
	GetBySlug(slug string) (*models.Organization, error)
	List(limit, offset int) ([]*models.Organization, error)
	Update(id int64, req *models.UpdateOrganizationRequest) error
	Delete(id int64) error
}

type organizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(req *models.CreateOrganizationRequest) (*models.Organization, error) {
	query := `
		INSERT INTO organizations (name, slug, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query, req.Name, req.Slug, req.Metadata, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	return r.GetByID(id)
}

func (r *organizationRepository) GetByID(id int64) (*models.Organization, error) {
	query := `
		SELECT id, name, slug, created_at, updated_at, is_active, metadata
		FROM organizations
		WHERE id = ?
	`
	org := &models.Organization{}
	err := r.db.QueryRow(query, id).Scan(
		&org.ID, &org.Name, &org.Slug, &org.CreatedAt,
		&org.UpdatedAt, &org.IsActive, &org.Metadata,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return org, nil
}

func (r *organizationRepository) GetBySlug(slug string) (*models.Organization, error) {
	query := `
		SELECT id, name, slug, created_at, updated_at, is_active, metadata
		FROM organizations
		WHERE slug = ?
	`
	org := &models.Organization{}
	err := r.db.QueryRow(query, slug).Scan(
		&org.ID, &org.Name, &org.Slug, &org.CreatedAt,
		&org.UpdatedAt, &org.IsActive, &org.Metadata,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return org, nil
}

func (r *organizationRepository) List(limit, offset int) ([]*models.Organization, error) {
	query := `
		SELECT id, name, slug, created_at, updated_at, is_active, metadata
		FROM organizations
		WHERE is_active = 1
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	orgs := []*models.Organization{}
	for rows.Next() {
		org := &models.Organization{}
		err := rows.Scan(
			&org.ID, &org.Name, &org.Slug, &org.CreatedAt,
			&org.UpdatedAt, &org.IsActive, &org.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (r *organizationRepository) Update(id int64, req *models.UpdateOrganizationRequest) error {

	query := "UPDATE organizations SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Name != nil {
		query += ", name = ?"
		args = append(args, *req.Name)
	}
	if req.IsActive != nil {
		query += ", is_active = ?"
		args = append(args, *req.IsActive)
	}
	if req.Metadata != nil {
		query += ", metadata = ?"
		args = append(args, *req.Metadata)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("organization not found")
	}

	return nil
}

func (r *organizationRepository) Delete(id int64) error {
	query := "UPDATE organizations SET is_active = 0 WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("organization not found")
	}

	return nil
}
