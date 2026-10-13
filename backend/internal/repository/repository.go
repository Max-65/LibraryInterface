package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Max-65/libraryInterface/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GroupByCode(ctx context.Context, code string) (models.Group, error) {
	const q = `
SELECT
	id, code, specialty_code, admission_year, institute_no, department_no, group_no,
	(
		EXTRACT(YEAR FROM CURRENT_DATE)
		- admission_year
		- CASE WHEN EXTRACT(MONTH FROM CURRENT_DATE) < 9 THEN 1 ELSE 0 END
		+ 1
	)::int
FROM groups
WHERE code = $1;
`
	var g models.Group
	err := r.db.QueryRow(ctx, q, code).Scan(
		&g.ID, &g.Code, &g.SpecialtyCode, &g.AdmissionYear,
		&g.InstituteNo, &g.DepartmentNo, &g.GroupNo, &g.CurrentCourse,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Group{}, ErrNotFound
		}
		return models.Group{}, err
	}
	return g, nil
}

func (r *Repository) GroupStudents(ctx context.Context, code string) ([]models.Student, error) {
	const q = `
SELECT 	s.id,
		s.firstname,
		s.lastname,
		s.groupid
FROM students s
JOIN groups g ON g.id = s.groupid
WHERE g.code = $1
ORDER BY s.lastname, s.firstname;
`
	rows, err := r.db.Query(ctx, q, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.ID, &s.Firstname, &s.Lastname, &s.GroupID); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

func (r *Repository) GroupRecommendations(ctx context.Context, code string) ([]models.Book, error) {
	const q = `
SELECT
    b.id,
    b.name,
    b.authors,
    b.disciplineid,
    b.course,
    b.isbn,
    b.publicationyear,
    b.publisher,
    COUNT(bc.id) FILTER (WHERE bc.status = 'в наличии')::bigint AS available_count,
    COUNT(bc.id)::bigint AS total_count
FROM groups g
JOIN disciplines d ON d.specialtyid = g.specialty_code
JOIN books b ON b.disciplineid = d.id
LEFT JOIN bookcopies bc ON bc.bookid = b.id
WHERE g.code = $1
  AND b.course = (
        EXTRACT(YEAR FROM CURRENT_DATE)
        - g.admission_year
        - CASE WHEN EXTRACT(MONTH FROM CURRENT_DATE) < 9 THEN 1 ELSE 0 END
        + 1
  )::int
GROUP BY
    b.id, b.name, b.authors, b.disciplineid, b.course,
    b.isbn, b.publicationyear, b.publisher
ORDER BY b.name;
`
	return r.scanBooks(ctx, q, code)
}

func (r *Repository) StudentLoans(ctx context.Context, studentID int64) ([]models.Loan, error) {
	const q = `
SELECT
	i.id, b.id, b.name, b.authors, i.copyid, i.issuedate
FROM issuing i
JOIN bookcopies bc ON bc.id = i.copyid
JOIN books b ON b.id = bc.bookid
WHERE i.studentid = $1 AND i.returndate IS NULL
ORDER BY i.issuedate DESC;
`
	rows, err := r.db.Query(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []models.Loan
	for rows.Next() {
		var l models.Loan
		if err := rows.Scan(&l.ID, &l.BookID, &l.BookName, &l.Authors, &l.CopyID, &l.IssueDate); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	return loans, rows.Err()
}

func (r *Repository) SearchBooks(ctx context.Context, q string) ([]models.Book, error) {
	const sqlQuery = `
SELECT
    b.id,
    b.name,
    b.authors,
    b.disciplineid,
    b.course,
    b.isbn,
    b.publicationyear,
    b.publisher,
    COUNT(bc.id) FILTER (WHERE bc.status = 'в наличии')::bigint AS available_count,
    COUNT(bc.id)::bigint AS total_count
FROM books b
LEFT JOIN bookcopies bc ON bc.bookid = b.id
WHERE b.name ILIKE '%' || $1 || '%'
GROUP BY
    b.id, b.name, b.authors, b.disciplineid, b.course,
    b.isbn, b.publicationyear, b.publisher
ORDER BY b.name;
`
	return r.scanBooks(ctx, sqlQuery, q)
}

func (r *Repository) BookByID(ctx context.Context, id int64) (models.Book, error) {
	const q = `
SELECT
	b.id, b.name, b.authors, b.disciplineid, b.course, b.isbn, b.publicationyear, b.publisher,
	COUNT(bc.id) FILTER (WHERE bc.status = 'в наличии')::bigint,
	COUNT(bc.id)::bigint
FROM books b
LEFT JOIN bookcopies bc ON bc.bookid = b.id
WHERE b.id = $1
GROUP BY b.id;
`
	var b models.Book
	err := r.db.QueryRow(ctx, q, id).Scan(
		&b.ID, &b.Name, &b.Authors, &b.DisciplineID, &b.Course,
		&b.ISBN, &b.PublicationYear, &b.Publisher,
		&b.AvailableCount, &b.TotalCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Book{}, ErrNotFound
		}
		return models.Book{}, err
	}
	return b, nil
}

func (r *Repository) BookCopies(ctx context.Context, bookID int64) ([]int64, error) {
	const q = `
SELECT id
FROM bookcopies
WHERE bookid = $1 AND status = 'в наличии'
ORDER BY id;
`
	rows, err := r.db.Query(ctx, q, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *Repository) CreateIssuing(ctx context.Context, req models.IssueRequest) (models.Issuing, error) {
	const q = `
INSERT INTO issuing (studentid, copyid, issuedate, returndate)
VALUES ($1, $2, CURRENT_DATE, NULL)
RETURNING id, studentid, copyid, issuedate, returndate;
`
	var out models.Issuing
	var ret sql.NullTime

	err := r.db.QueryRow(ctx, q, req.StudentID, req.CopyID).Scan(
		&out.ID, &out.StudentID, &out.CopyID, &out.IssueDate, &ret,
	)
	if err != nil {
		return models.Issuing{}, err
	}
	if ret.Valid {
		out.ReturnDate = &ret.Time
	}
	return out, nil
}

func (r *Repository) ReturnIssuing(ctx context.Context, id int64) (models.Issuing, error) {
	const q = `
UPDATE issuing
SET returndate = CURRENT_DATE
WHERE id = $1 AND returndate IS NULL
RETURNING id, studentid, copyid, issuedate, returndate;
`
	var out models.Issuing
	var ret sql.NullTime

	err := r.db.QueryRow(ctx, q, id).Scan(
		&out.ID, &out.StudentID, &out.CopyID, &out.IssueDate, &ret,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Issuing{}, ErrNotFound
		}
		return models.Issuing{}, err
	}
	if ret.Valid {
		out.ReturnDate = &ret.Time
	}
	return out, nil
}

func (r *Repository) scanBooks(ctx context.Context, q string, args ...any) ([]models.Book, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(
			&b.ID, &b.Name, &b.Authors, &b.DisciplineID, &b.Course,
			&b.ISBN, &b.PublicationYear, &b.Publisher, &b.AvailableCount, &b.TotalCount,
		); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func asNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
