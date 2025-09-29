package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Difficulty string

const (
	Beginner     Difficulty = "Beginner"
	Intermediate Difficulty = "Intermediate"
	Advanced     Difficulty = "Advanced"
)

func (d Difficulty) Valid() bool {
	return d == Beginner || d == Intermediate || d == Advanced
}

type Course struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Difficulty  Difficulty `json:"difficulty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Enrollment struct {
	ID           uuid.UUID `json:"id"`
	StudentEmail string    `json:"student_email"`
	CourseID     uuid.UUID `json:"course_id"`
	EnrolledAt   time.Time `json:"enrolled_at"`
}

var (
	ErrTitleEmpty   = errors.New("title must not be empty")
	ErrDifficulty   = errors.New("invalid difficulty")
	ErrEmailInvalid = errors.New("invalid email format")
	emailRe         = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
)

func ValidateNewCourse(title string, diff Difficulty) error {
	if strings.TrimSpace(title) == "" {
		return ErrTitleEmpty
	}
	if !diff.Valid() {
		return ErrDifficulty
	}
	return nil
}

func ValidateEmail(s string) error {
	if !emailRe.MatchString(s) {
		return ErrEmailInvalid
	}
	return nil
}
