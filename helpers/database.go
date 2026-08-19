package helpers

import (
	"backend-test/internal/models"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupPostgreSQL() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", GetEnv("DB_HOST", "127.0.0.1"), GetEnv("DB_PORT", "5432"), GetEnv("DB_USER", "postgres"), GetEnv("DB_PASSWORD", ""), GetEnv("DB_NAME", "maintenance_request_log"), GetEnv("DB_SSLMODE", "disable"))
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	DB = database
	if err := DB.AutoMigrate(&models.User{}, &models.MaintenanceRequest{}, &models.RequestActivity{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	if err := SeedDatabase(); err != nil {
		log.Fatal("failed to seed database: ", err)
	}
	Logger.Info("database connected, migrated, and seeded")
}

func SeedDatabase() error {
	password, err := bcrypt.GenerateFromPassword([]byte(GetEnv("SEED_PASSWORD", "Password123!")), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	seedUsers := []models.User{
		{ID: 1, Name: "Budi Santoso", Email: "admin@industrialops.com", PasswordHash: string(password), Role: models.RoleAdmin, Status: models.UserActive},
		{ID: 2, Name: "Andi Pratama", Email: "supervisor@industrialops.com", PasswordHash: string(password), Role: models.RoleSupervisor, Status: models.UserActive},
		{ID: 3, Name: "Ahmad Fauzi", Email: "operator@industrialops.com", PasswordHash: string(password), Role: models.RoleOperator, Status: models.UserActive},
		{ID: 4, Name: "Sarah Jenkins", Email: "sarah.j@industrialops.com", PasswordHash: string(password), Role: models.RoleOperator, Status: models.UserActive},
		{ID: 5, Name: "Michael Chen", Email: "michael.c@industrialops.com", PasswordHash: string(password), Role: models.RoleOperator, Status: models.UserInactive},
	}
	for i := range seedUsers {
		var user models.User
		if err := DB.Where("email = ?", seedUsers[i].Email).Attrs(seedUsers[i]).FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}
	if err := DB.Exec("SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 1), true) FROM users").Error; err != nil {
		return err
	}
	var users []models.User
	if err := DB.Order("id ASC").Find(&users).Error; err != nil {
		return err
	}
	byEmail := map[string]models.User{}
	for _, user := range users {
		byEmail[user.Email] = user
	}
	admin, supervisor, operator := byEmail["admin@industrialops.com"], byEmail["supervisor@industrialops.com"], byEmail["operator@industrialops.com"]
	technician := "Sarah Jenkins (Hydraulics Spec.)"
	target := time.Now().AddDate(0, 0, 7)
	reviewed := time.Now().Add(-2 * time.Hour)
	seedRequests := []models.MaintenanceRequest{
		{ID: "MR-0001", MachineID: "CNC-Lathe-04", Problem: "Spindle vibration exceeding tolerance during high-speed operation.", Priority: models.PriorityCritical, Status: models.StatusSubmitted, CreatedByID: operator.ID, CreatedAt: time.Now().Add(-4 * time.Hour)},
		{ID: "MR-0002", MachineID: "PRESS-03", Problem: "Hydraulic pressure drops after 30 minutes of operation. Machine triggers auto-shutoff safety protocol. Requires immediate inspection of main seal and fluid lines.", Priority: models.PriorityCritical, Status: models.StatusApproved, CreatedByID: admin.ID, ReviewedByID: &supervisor.ID, ReviewedAt: &reviewed, Technician: &technician, TargetDate: &target, CreatedAt: time.Now().Add(-26 * time.Hour)},
		{ID: "MR-0003", MachineID: "HVAC-Unit-Roof", Problem: "Filter replacement overdue and restricted airflow detected by monitoring system.", Priority: models.PriorityMedium, Status: models.StatusRejected, CreatedByID: operator.ID, ReviewedByID: &supervisor.ID, ReviewedAt: &reviewed, CreatedAt: time.Now().Add(-48 * time.Hour)},
		{ID: "MR-0004", MachineID: "Conveyor-B-12", Problem: "Squeaking noise from drive belt requires inspection before the next production cycle.", Priority: models.PriorityLow, Status: models.StatusApproved, CreatedByID: operator.ID, ReviewedByID: &supervisor.ID, ReviewedAt: &reviewed, CreatedAt: time.Now().Add(-72 * time.Hour)},
		{ID: "MR-0005", MachineID: "BOILER-02", Problem: "Temperature sensor readings fluctuate intermittently outside normal tolerance.", Priority: models.PriorityHigh, Status: models.StatusSubmitted, CreatedByID: admin.ID, CreatedAt: time.Now().Add(-96 * time.Hour)},
		{ID: "MR-0006", MachineID: "PACK-LINE-07", Problem: "Emergency stop button casing is loose and requires immediate replacement.", Priority: models.PriorityCritical, Status: models.StatusSubmitted, CreatedByID: operator.ID, CreatedAt: time.Now().Add(-120 * time.Hour)},
	}
	for i := range seedRequests {
		var request models.MaintenanceRequest
		result := DB.First(&request, "id = ?", seedRequests[i].ID)
		if result.Error == gorm.ErrRecordNotFound {
			if err := DB.Create(&seedRequests[i]).Error; err != nil {
				return err
			}
			request = seedRequests[i]
		} else if result.Error != nil {
			return result.Error
		}
		var count int64
		if err := DB.Model(&models.RequestActivity{}).Where("request_id = ?", request.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			creator := admin
			if request.CreatedByID == operator.ID {
				creator = operator
			}
			created := models.RequestActivity{RequestID: request.ID, ActorID: creator.ID, Type: models.ActivityCreated, Title: "Request Created", Description: fmt.Sprintf("Submitted by %s via web portal.", creator.Name), CreatedAt: request.CreatedAt}
			if err := DB.Create(&created).Error; err != nil {
				return err
			}
			if request.Status != models.StatusSubmitted {
				activityType, title := models.ActivityApproved, "Request Approved"
				if request.Status == models.StatusRejected {
					activityType, title = models.ActivityRejected, "Request Rejected"
				}
				activity := models.RequestActivity{RequestID: request.ID, ActorID: supervisor.ID, Type: activityType, Title: title, Description: fmt.Sprintf("%s by %s.", request.Status, supervisor.Name), CreatedAt: reviewed}
				if err := DB.Create(&activity).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}
