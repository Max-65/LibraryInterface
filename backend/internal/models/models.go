package models

import "time"

type Group struct {
	ID            int64  `json:"id"`
	Code          string `json:"code"`
	SpecialtyCode string `json:"specialty_code"`
	AdmissionYear int    `json:"admission_year"`
	InstituteNo   int    `json:"institute_no"`
	DepartmentNo  int    `json:"department_no"`
	GroupNo       int    `json:"group_no"`
	CurrentCourse int    `json:"current_course"`
}

type Student struct {
	ID        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	GroupID   int64  `json:"group_id"`
}

type Loan struct {
	ID        int64     `json:"id"`
	BookID    int64     `json:"book_id"`
	BookName  string    `json:"book_name"`
	Authors   []string  `json:"authors"`
	CopyID    int64     `json:"copy_id"`
	IssueDate time.Time `json:"issue_date"`
}

type Book struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Authors         []string `json:"authors"`
	DisciplineID    int64    `json:"disciplineid"`
	Course          int      `json:"course"`
	ISBN            string   `json:"isbn"`
	PublicationYear int      `json:"publicationyear"`
	Publisher       string   `json:"publisher"`
	AvailableCount  int64    `json:"available_count,omitempty"`
	TotalCount      int64    `json:"total_count,omitempty"`
}

type Issuing struct {
	ID         int64      `json:"id"`
	StudentID  int64      `json:"student_id"`
	CopyID     int64      `json:"copy_id"`
	IssueDate  time.Time  `json:"issuedate"`
	ReturnDate *time.Time `json:"returndate,omitempty"`
}

type IssueRequest struct {
	StudentID int64 `json:"student_id"`
	CopyID    int64 `json:"copy_id"`
}
