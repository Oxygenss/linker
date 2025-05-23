package repository

import (
	"fmt"

	"github.com/Oxygenss/linker/internal/config"
	"github.com/Oxygenss/linker/internal/models"
	"github.com/Oxygenss/linker/internal/repository/postgres"
	"github.com/Oxygenss/linker/pkg/logger"
	"github.com/google/uuid"
)

type Repository struct {
	logger            logger.Logger
	StudentRepository StudentRepository
	TeacherRepository TeacherRepository
	UserRepository    UserRepository
	RequestRepository RequestRepository
	WorkRepository    WorkRepository
}

type StudentRepository interface {
	GetByTelegramID(telegramID int64) (models.Student, error)
	GetByID(id string) (models.Student, error)
	Create(student models.Student) (uuid.UUID, error)
	Update(student models.Student) error
	GetAll() ([]models.Student, error)
	Search(search string) ([]models.Student, error)
	Delete(id string) error
}

type TeacherRepository interface {
	GetByTelegramID(telegramID int64) (models.Teacher, error)
	GetByID(id string) (models.Teacher, error)
	Update(teacher models.Teacher) error
	Create(teacher models.Teacher) (uuid.UUID, error)
	GetAll() ([]models.Teacher, error)
	Search(search string) ([]models.Teacher, error)
	Delete(id string) error
}

type UserRepository interface {
	GetRoleByID(id string) (string, error)
	GetRoleByTelegramID(telegramID int64) (string, error)
}

type RequestRepository interface {
	Create(models.Request) error
}

type WorkRepository interface {
	Create(models.Work) error
	GetAll(userID uuid.UUID) ([]models.Work, error)
	Delete(id uuid.UUID) error
}

func NewRepository(config *config.Config, logger *logger.Logger) (*Repository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	db, err := postgres.NewPostgresConnection(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	postgresRepository := postgres.NewPostgresRepository(db, logger)

	return &Repository{
		StudentRepository: postgresRepository.StudentRepository,
		TeacherRepository: postgresRepository.TeacherRepository,
		UserRepository:    postgresRepository.UserRepository,
		RequestRepository: postgresRepository.RequestRepository,
		WorkRepository:    postgresRepository.WorksRepository,
	}, nil
}
