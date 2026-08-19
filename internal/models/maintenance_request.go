package models

import "time"

type RequestPriority string
type RequestStatus string

type ActivityType string

const (
	PriorityLow      RequestPriority = "Low"
	PriorityMedium   RequestPriority = "Medium"
	PriorityHigh     RequestPriority = "High"
	PriorityCritical RequestPriority = "Critical"

	StatusSubmitted RequestStatus = "Submitted"
	StatusApproved  RequestStatus = "Approved"
	StatusRejected  RequestStatus = "Rejected"

	ActivityCreated  ActivityType = "created"
	ActivityUpdated  ActivityType = "updated"
	ActivityApproved ActivityType = "approved"
	ActivityRejected ActivityType = "rejected"
)

type MaintenanceRequest struct {
	ID           string            `gorm:"primaryKey;size:20" json:"id"`
	MachineID    string            `gorm:"size:100;not null;index" json:"machineId"`
	Problem      string            `gorm:"type:text;not null" json:"problem"`
	Priority     RequestPriority   `gorm:"type:varchar(20);not null;index" json:"priority"`
	Status       RequestStatus     `gorm:"type:varchar(20);not null;default:Submitted;index" json:"status"`
	CreatedByID  uint              `gorm:"not null;index" json:"createdById"`
	CreatedBy    User              `gorm:"foreignKey:CreatedByID" json:"createdBy"`
	ReviewedByID *uint             `gorm:"index" json:"reviewedById,omitempty"`
	ReviewedBy   *User             `gorm:"foreignKey:ReviewedByID" json:"reviewedBy,omitempty"`
	ReviewedAt   *time.Time        `json:"reviewedAt,omitempty"`
	Technician   *string           `gorm:"size:180" json:"technician,omitempty"`
	TargetDate   *time.Time        `gorm:"type:date" json:"targetDate,omitempty"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	Activities   []RequestActivity `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE" json:"activities,omitempty"`
}

type RequestActivity struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	RequestID   string       `gorm:"size:20;not null;index" json:"requestId"`
	ActorID     uint         `gorm:"not null" json:"actorId"`
	Actor       User         `gorm:"foreignKey:ActorID" json:"actor"`
	Type        ActivityType `gorm:"type:varchar(20);not null" json:"type"`
	Title       string       `gorm:"size:160;not null" json:"title"`
	Description string       `gorm:"type:text;not null" json:"description"`
	CreatedAt   time.Time    `json:"createdAt"`
}
