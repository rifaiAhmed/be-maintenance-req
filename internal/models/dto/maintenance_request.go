package dto

import "backend-test/internal/models"

type CreateMaintenanceRequest struct {
	MachineID string                 `json:"machineId" binding:"required,min=2,max=100"`
	Problem   string                 `json:"problem" binding:"required,min=10,max=500"`
	Priority  models.RequestPriority `json:"priority" binding:"required,oneof=Low Medium High Critical"`
}

type UpdateMaintenanceRequest struct {
	MachineID  string                 `json:"machineId" binding:"omitempty,min=2,max=100"`
	Problem    string                 `json:"problem" binding:"required,min=10,max=500"`
	Priority   models.RequestPriority `json:"priority" binding:"omitempty,oneof=Low Medium High Critical"`
	Technician *string                `json:"technician" binding:"omitempty,max=180"`
	TargetDate *string                `json:"targetDate" binding:"omitempty,datetime=2006-01-02"`
}

type ReviewMaintenanceRequest struct {
	Status models.RequestStatus `json:"status" binding:"required,oneof=Approved Rejected"`
}

type RequestFilter struct {
	Search   string
	Status   string
	Priority string
	Page     int
	Limit    int
}
