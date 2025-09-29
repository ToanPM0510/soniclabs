package pg

import (
	"context"
	"errors"
	"strings"
	"time"

	"toanpm0510/soniclabs/internal/domain"

	"github.com/google/uuid"
)

var (
	ErrDuplicateEnrollment = errors.New("duplicate enrollment")
	ErrCourseNotFound      = errors.New("course not found")
)

func CreateCourse(ctx context.Context, db *DB, title string, desc *string, diff domain.Difficulty) (domain.Course, error) {
	id := uuid.New()
	var createdAt time.Time
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO courses (id, title, description, difficulty)
		VALUES ($1,$2,$3,$4)
		RETURNING created_at`,
		id, title, desc, string(diff),
	).Scan(&createdAt)
	if err != nil {
		return domain.Course{}, err
	}
	return domain.Course{
		ID:          id,
		Title:       title,
		Description: desc,
		Difficulty:  diff,
		CreatedAt:   createdAt,
	}, nil
}

func ListCourses(ctx context.Context, db *DB) ([]domain.Course, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, title, description, difficulty, created_at
		FROM courses
		ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Course
	for rows.Next() {
		var c domain.Course
		var diff string
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &diff, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		c.Difficulty = domain.Difficulty(diff)
		out = append(out, c)
	}
	return out, rows.Err()
}
func Enroll(ctx context.Context, db *DB, email string, courseID uuid.UUID) (domain.Enrollment, error) {
	id := uuid.New()
	var returnedID uuid.UUID
	var enrolledAt time.Time

	err := db.Pool.QueryRow(ctx, `
        INSERT INTO enrollments (id, student_email, course_id)
        VALUES ($1, lower($2), $3)
        ON CONFLICT (lower(student_email), course_id) 
        DO UPDATE SET student_email = EXCLUDED.student_email
        RETURNING id, enrolled_at
    `, id, email, courseID).Scan(&returnedID, &enrolledAt)

	if err != nil {
		// Check for foreign key violation (course doesn't exist)
		if strings.Contains(err.Error(), "foreign key constraint") {
			return domain.Enrollment{}, ErrCourseNotFound
		}
		return domain.Enrollment{}, err
	}

	// If returned ID != new ID, enrollment already existed
	if returnedID != id {
		return domain.Enrollment{}, ErrDuplicateEnrollment
	}

	return domain.Enrollment{
		ID:           returnedID,
		StudentEmail: strings.ToLower(email),
		CourseID:     courseID,
		EnrolledAt:   enrolledAt,
	}, nil
}

func ListEnrollmentsByEmail(ctx context.Context, db *DB, email string) ([]domain.Enrollment, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, student_email, course_id, enrolled_at
		FROM enrollments
		WHERE student_email = $1
		ORDER BY enrolled_at DESC, id DESC
	`, strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Enrollment
	for rows.Next() {
		var e domain.Enrollment
		if err := rows.Scan(&e.ID, &e.StudentEmail, &e.CourseID, &e.EnrolledAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
