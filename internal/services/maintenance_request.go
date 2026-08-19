package service

import (
	"backend-test/helpers"
	"backend-test/internal/interfaces"
	"backend-test/internal/models"
	"backend-test/internal/models/dto"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type RequestService struct{ Requests interfaces.RequestRepository }

func NewRequestService(requests interfaces.RequestRepository) *RequestService {
	return &RequestService{Requests: requests}
}
func canView(actor *models.User, request *models.MaintenanceRequest) bool {
	return actor.Role != models.RoleOperator || request.CreatedByID == actor.ID
}

func (s *RequestService) List(ctx context.Context, actor *models.User, filter dto.RequestFilter) ([]models.MaintenanceRequest, models.Pagination, error) {
	requests, total, err := s.Requests.List(ctx, filter, actor)
	return requests, models.Pagination{Page: filter.Page, Limit: filter.Limit, Total: total, TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit)))}, err
}
func (s *RequestService) Get(ctx context.Context, actor *models.User, id string) (*models.MaintenanceRequest, error) {
	request, err := s.Requests.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, helpers.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if !canView(actor, request) {
		return nil, helpers.ErrForbidden
	}
	return request, nil
}
func (s *RequestService) Create(ctx context.Context, actor *models.User, req dto.CreateMaintenanceRequest) (*models.MaintenanceRequest, error) {
	request := &models.MaintenanceRequest{MachineID: req.MachineID, Problem: req.Problem, Priority: req.Priority, Status: models.StatusSubmitted, CreatedByID: actor.ID}
	activity := &models.RequestActivity{ActorID: actor.ID, Type: models.ActivityCreated, Title: "Request Created", Description: fmt.Sprintf("Submitted by %s via web portal.", actor.Name)}
	if err := s.Requests.Create(ctx, request, activity); err != nil {
		return nil, err
	}
	return s.Requests.FindByID(ctx, request.ID)
}
func (s *RequestService) Update(ctx context.Context, actor *models.User, id string, req dto.UpdateMaintenanceRequest) (*models.MaintenanceRequest, error) {
	request, err := s.Get(ctx, actor, id)
	if err != nil {
		return nil, err
	}
	if actor.Role != models.RoleAdmin && (request.CreatedByID != actor.ID || request.Status != models.StatusSubmitted) {
		return nil, helpers.ErrForbidden
	}
	if req.MachineID != "" {
		request.MachineID = req.MachineID
	}
	if req.Priority != "" {
		request.Priority = req.Priority
	}
	request.Problem = req.Problem
	request.Technician = req.Technician
	if req.TargetDate != nil {
		date, parseErr := time.Parse("2006-01-02", *req.TargetDate)
		if parseErr != nil {
			return nil, helpers.ErrInvalid
		}
		request.TargetDate = &date
	} else {
		request.TargetDate = nil
	}
	activity := &models.RequestActivity{ActorID: actor.ID, Type: models.ActivityUpdated, Title: "Request Updated", Description: fmt.Sprintf("Request details updated by %s.", actor.Name)}
	if err := s.Requests.Update(ctx, request, activity); err != nil {
		return nil, err
	}
	return s.Requests.FindByID(ctx, id)
}
func (s *RequestService) Review(ctx context.Context, actor *models.User, id string, req dto.ReviewMaintenanceRequest) (*models.MaintenanceRequest, error) {
	if actor.Role != models.RoleSupervisor && actor.Role != models.RoleAdmin {
		return nil, helpers.ErrForbidden
	}
	request, err := s.Requests.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, helpers.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	now := time.Now()
	request.Status, request.ReviewedByID, request.ReviewedAt = req.Status, &actor.ID, &now
	typeName, title := models.ActivityApproved, "Request Approved"
	if req.Status == models.StatusRejected {
		typeName, title = models.ActivityRejected, "Request Rejected"
	}
	activity := &models.RequestActivity{ActorID: actor.ID, Type: typeName, Title: title, Description: fmt.Sprintf("%s by %s.", req.Status, actor.Name)}
	if err := s.Requests.Review(ctx, request, activity); err != nil {
		return nil, err
	}
	return s.Requests.FindByID(ctx, id)
}
func (s *RequestService) Delete(ctx context.Context, actor *models.User, id string) error {
	if actor.Role != models.RoleAdmin {
		return helpers.ErrForbidden
	}
	if _, err := s.Requests.FindByID(ctx, id); errors.Is(err, gorm.ErrRecordNotFound) {
		return helpers.ErrNotFound
	} else if err != nil {
		return err
	}
	return s.Requests.Delete(ctx, id)
}
func (s *RequestService) Dashboard(ctx context.Context, actor *models.User) (map[string]int64, []models.MaintenanceRequest, error) {
	counts, err := s.Requests.DashboardCounts(ctx, actor)
	if err != nil {
		return nil, nil, err
	}
	items, _, err := s.Requests.List(ctx, dto.RequestFilter{Page: 1, Limit: 5}, actor)
	return counts, items, err
}
