package gui

type ViewModel struct {
	AppName string
	Title   string
	Active  string
	Notice  string
	Error   string

	Patients      []PatientRow
	Appointments  []AppointmentRow
	Slots         []ScheduleSlotRow
	FreeSlots     []ScheduleSlotRow
	Visits        []VisitRow
	Doctors       []DoctorRow
	DoctorOptions []DoctorRow

	PatientID     string
	PatientSearch string
	SearchEmail   string
	DoctorID      string
	DoctorSearch  string
	AppointmentID string
	Date          string
	DateTime      string
}

type PatientRow struct {
	ID        string
	FullName  string
	Phone     string
	Email     string
	BirthDate string
	CreatedAt string
}

type AppointmentRow struct {
	ID              string
	PatientID       string
	PatientName     string
	PatientEmail    string
	DoctorID        string
	DoctorName      string
	DoctorSpecialty string
	Status          string
	Date            string
	Badge           string
}

type ScheduleSlotRow struct {
	ID              string
	DoctorID        string
	DoctorName      string
	DoctorSpecialty string
	Date            string
	StartTime       string
	EndTime         string
	Cabinet         string
	IsBooked        bool
	Badge           string
}

type VisitRow struct {
	ID              string
	AppointmentID   string
	Complaints      string
	Diagnosis       string
	Recommendations string
	Status          string
	Badge           string
}

type DoctorRow struct {
	ID        string
	FullName  string
	Specialty string
	Cabinet   string
}
