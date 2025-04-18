package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Oxygenss/linker/internal/models"
	"github.com/Oxygenss/linker/pkg/logger"
	"github.com/google/uuid"
)

type TeacherRepository struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewTeacherRepository(db *sql.DB, logger *logger.Logger) *TeacherRepository {
	return &TeacherRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TeacherRepository) GetByID(id string) (models.Teacher, error) {
	r.logger.Info("[TeacherRepository: GetByID]", "id", id)

	if r.db == nil {
		return models.Teacher{}, fmt.Errorf("database connection is not initialized")
	}

	_, err := uuid.Parse(id)
	if err != nil {
		return models.Teacher{}, fmt.Errorf("invalid UUID format: %w", err)
	}

	query := `SELECT id, telegram_id, first_name, middle_name, last_name, 
                     degree, position, department, is_free, idea, about 
              FROM teachers WHERE id = $1`

	var teacher models.Teacher
	err = r.db.QueryRow(query, id).Scan(
		&teacher.ID,
		&teacher.TelegramID,
		&teacher.FirstName,
		&teacher.MiddleName,
		&teacher.LastName,
		&teacher.Degree,
		&teacher.Position,
		&teacher.Department,
		&teacher.IsFree,
		&teacher.Idea,
		&teacher.About,
	)

	switch {
	case err == nil:
		return teacher, nil
	case errors.Is(err, sql.ErrNoRows):
		return models.Teacher{}, fmt.Errorf("teacher not found: %w", err)
	default:
		r.logger.Error("Failed to get teacher by ID", "error", err, "id", id)
		return models.Teacher{}, fmt.Errorf("failed to get teacher: %w", err)
	}
}

func (r *TeacherRepository) GetByTelegramID(telegramID int64) (models.Teacher, error) {
	r.logger.Info("[TeacherRepository: GetByTelegramID]", "telegramID", telegramID)

	if r.db == nil {
		return models.Teacher{}, fmt.Errorf("database connection is not initialized")
	}

	query := `SELECT id, telegram_id, first_name, middle_name, last_name,
                     degree, position, department, is_free, idea, about
              FROM teachers WHERE telegram_id = $1`

	var teacher models.Teacher
	err := r.db.QueryRow(query, telegramID).Scan(
		&teacher.ID,
		&teacher.TelegramID,
		&teacher.FirstName,
		&teacher.MiddleName,
		&teacher.LastName,
		&teacher.Degree,
		&teacher.Position,
		&teacher.Department,
		&teacher.IsFree,
		&teacher.Idea,
		&teacher.About,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Teacher{}, fmt.Errorf("teacher not found")
		}
		r.logger.Error("Failed to get teacher by Telegram ID", "error", err, "telegramID", telegramID)
		return models.Teacher{}, fmt.Errorf("failed to retrieve teacher: %w", err)
	}

	return teacher, nil
}

func (r *TeacherRepository) GetAll() ([]models.Teacher, error) {
	r.logger.Info("[TeacherRepository: GetAll]")

	query := `SELECT id, telegram_id, first_name, middle_name, last_name,
                     degree, position, department, is_free, idea, about
              FROM teachers`

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to retrieve teachers", "error", err)
		return nil, fmt.Errorf("failed to retrieve teachers: %w", err)
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(
			&teacher.ID,
			&teacher.TelegramID,
			&teacher.FirstName,
			&teacher.MiddleName,
			&teacher.LastName,
			&teacher.Degree,
			&teacher.Position,
			&teacher.Department,
			&teacher.IsFree,
			&teacher.Idea,
			&teacher.About,
		)
		if err != nil {
			r.logger.Error("Failed to scan teacher", "error", err)
			return nil, fmt.Errorf("failed to scan teacher: %w", err)
		}
		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Error after scanning teachers", "error", err)
		return nil, fmt.Errorf("error after scanning teachers: %w", err)
	}

	return teachers, nil
}

func (r *TeacherRepository) Create(teacher models.Teacher) (uuid.UUID, error) {
	r.logger.Info("[TeacherRepository: Create]")

	if r.db == nil {
		return uuid.Nil, fmt.Errorf("database connection is not initialized")
	}

	teacher.ID = uuid.New()

	query := `INSERT INTO teachers 
              (id, telegram_id, first_name, middle_name, last_name, 
               degree, position, department, is_free, idea, about) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(query,
		teacher.ID,
		teacher.TelegramID,
		teacher.FirstName,
		teacher.MiddleName,
		teacher.LastName,
		teacher.Degree,
		teacher.Position,
		teacher.Department,
		teacher.IsFree,
		teacher.Idea,
		teacher.About,
	)

	if err != nil {
		r.logger.Error("Failed to create teacher", "error", err, "teacherID", teacher.ID)
		return uuid.Nil, fmt.Errorf("failed to insert teacher: %w", err)
	}

	return teacher.ID, nil
}

func (r *TeacherRepository) Update(teacher models.Teacher) error {
	r.logger.Info("[TeacherRepository: Update]", "teacherID", teacher.ID)

	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	query := `
        UPDATE teachers 
        SET 
            first_name = $1,
            middle_name = $2,
            last_name = $3,
            degree = $4,
            position = $5,
            department = $6,
            is_free = $7,
            idea = $8,
            about = $9
        WHERE id = $10
    `

	result, err := r.db.Exec(query,
		teacher.FirstName,
		teacher.MiddleName,
		teacher.LastName,
		teacher.Degree,
		teacher.Position,
		teacher.Department,
		teacher.IsFree,
		teacher.Idea,
		teacher.About,
		teacher.ID,
	)

	if err != nil {
		r.logger.Error("Failed to update teacher", "error", err, "teacherID", teacher.ID)
		return fmt.Errorf("failed to update teacher: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Warn("Teacher not found for update", "teacherID", teacher.ID)
		return fmt.Errorf("teacher with ID %s not found", teacher.ID)
	}

	return nil
}
