package repository

import (
	"backend-test/internal/models"
	"backend-test/internal/models/dto"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type RequestRepository struct{ Database *gorm.DB }

func (r *RequestRepository) DB() *gorm.DB { return r.Database }

func requestPreloads(db *gorm.DB) *gorm.DB {
	return db.Preload("CreatedBy").Preload("ReviewedBy").Preload("Activities", func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") }).Preload("Activities.Actor")
}

func (r *RequestRepository) List(ctx context.Context, filter dto.RequestFilter, viewer *models.User) ([]models.MaintenanceRequest, int64, error) {
	var requests []models.MaintenanceRequest
	var total int64
	query := r.Database.WithContext(ctx).Model(&models.MaintenanceRequest{})
	if viewer.Role == models.RoleOperator {
		query = query.Where("created_by_id = ?", viewer.ID)
	}
	if filter.Search != "" {
		term := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(id) LIKE ? OR LOWER(machine_id) LIKE ? OR LOWER(problem) LIKE ?", term, term, term)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query = requestPreloads(query)
	if err := query.Order("created_at DESC").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

func (r *RequestRepository) FindByID(ctx context.Context, id string) (*models.MaintenanceRequest, error) {
	var request models.MaintenanceRequest
	if err := requestPreloads(r.Database.WithContext(ctx)).First(&request, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *RequestRepository) Create(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error {
	return r.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", 721904).Error; err != nil {
			return err
		}
		var maxNumber int
		if err := tx.Raw("SELECT COALESCE(MAX(CAST(SUBSTRING(id FROM 4) AS INTEGER)), 0) FROM maintenance_requests").Scan(&maxNumber).Error; err != nil {
			return err
		}
		request.ID = fmt.Sprintf("MR-%04d", maxNumber+1)
		if err := tx.Create(request).Error; err != nil {
			return err
		}
		activity.RequestID = request.ID
		return tx.Create(activity).Error
	})
}

func (r *RequestRepository) Update(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error {
	return r.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(request).Error; err != nil {
			return err
		}
		activity.RequestID = request.ID
		return tx.Create(activity).Error
	})
}
func (r *RequestRepository) Review(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error {
	return r.Update(ctx, request, activity)
}
func (r *RequestRepository) Delete(ctx context.Context, id string) error {
	return r.Database.WithContext(ctx).Delete(&models.MaintenanceRequest{}, "id = ?", id).Error
}

func (r *RequestRepository) DashboardCounts(ctx context.Context, viewer *models.User) (map[string]int64, error) {
	counts := map[string]int64{"total": 0, "submitted": 0, "approved": 0, "rejected": 0}
	base := r.Database.WithContext(ctx).Model(&models.MaintenanceRequest{})
	if viewer.Role == models.RoleOperator {
		base = base.Where("created_by_id = ?", viewer.ID)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	counts["total"] = total
	for _, status := range []string{"Submitted", "Approved", "Rejected"} {
		var count int64
		query := r.Database.WithContext(ctx).Model(&models.MaintenanceRequest{}).Where("status = ?", status)
		if viewer.Role == models.RoleOperator {
			query = query.Where("created_by_id = ?", viewer.ID)
		}
		if err := query.Count(&count).Error; err != nil {
			return nil, err
		}
		counts[strings.ToLower(status)] = count
	}
	return counts, nil
}
