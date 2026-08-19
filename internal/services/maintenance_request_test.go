package service

import (
	"backend-test/helpers"
	"backend-test/internal/models"
	"backend-test/internal/models/dto"
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type fakeRequestRepo struct {
	req *models.MaintenanceRequest
	err error
}

func (f *fakeRequestRepo) List(ctx context.Context, filter dto.RequestFilter, viewer *models.User) ([]models.MaintenanceRequest, int64, error) {
	return nil, 0, nil
}

func (f *fakeRequestRepo) FindByID(ctx context.Context, id string) (*models.MaintenanceRequest, error) {
	return f.req, f.err
}

func (f *fakeRequestRepo) Create(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error {
	return nil
}

func (f *fakeRequestRepo) Update(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error {
	return nil
}

func (f *fakeRequestRepo) Review(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error {
	return nil
}

func (f *fakeRequestRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (f *fakeRequestRepo) DashboardCounts(ctx context.Context, viewer *models.User) (map[string]int64, error) {
	return nil, nil
}

func (f *fakeRequestRepo) DB() *gorm.DB {
	return nil
}

func TestRequestService_Get(t *testing.T) {
	ctx := context.Background()
	sample := &models.MaintenanceRequest{ID: "MR-001", CreatedByID: 10}

	tests := []struct {
		name    string
		actor   *models.User
		repoErr error
		wantErr error
		wantReq *models.MaintenanceRequest
	}{
		{
			name:    "admin can view any request",
			actor:   &models.User{ID: 1, Role: models.RoleAdmin},
			repoErr: nil,
			wantErr: nil,
			wantReq: sample,
		},
		{
			name:    "operator can view own request",
			actor:   &models.User{ID: 10, Role: models.RoleOperator},
			repoErr: nil,
			wantErr: nil,
			wantReq: sample,
		},
		{
			name:    "operator cannot view others request",
			actor:   &models.User{ID: 11, Role: models.RoleOperator},
			repoErr: nil,
			wantErr: helpers.ErrForbidden,
			wantReq: nil,
		},
		{
			name:    "not found returns ErrNotFound",
			actor:   &models.User{ID: 1, Role: models.RoleAdmin},
			repoErr: gorm.ErrRecordNotFound,
			wantErr: helpers.ErrNotFound,
			wantReq: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRequestRepo{req: sample, err: tt.repoErr}
			svc := NewRequestService(repo)
			got, err := svc.Get(ctx, tt.actor, "MR-001")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantReq != nil {
				if got == nil {
					t.Fatalf("Get() returned nil request, want %v", tt.wantReq)
				}
				if got.ID != tt.wantReq.ID {
					t.Fatalf("Get() request ID = %s, want %s", got.ID, tt.wantReq.ID)
				}
			}
		})
	}
}
