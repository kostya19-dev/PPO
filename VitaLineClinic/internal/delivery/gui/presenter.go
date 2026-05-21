package gui

import (
	"time"

	"vitaline-clinic/internal/domain"
	"vitaline-clinic/internal/usecase"
)

const (
	dateLayout     = "2006-01-02"
	timeLayout     = "15:04"
	dateTimeLayout = "2006-01-02 15:04"
	formTimeLayout = "2006-01-02T15:04"
)

type Presenter struct{}

func NewPresenter() *Presenter {
	return &Presenter{}
}

func (p *Presenter) Patients(items []usecase.PatientResponseDTO) []PatientRow {
	rows := make([]PatientRow, 0, len(items))
	for _, item := range items {
		rows = append(rows, PatientRow{
			ID:        item.ID.String(),
			FullName:  item.FullName,
			Phone:     item.Phone,
			Email:     item.Email,
			BirthDate: item.BirthDate.Format(dateLayout),
			CreatedAt: item.CreatedAt.Format(dateTimeLayout),
		})
	}
	return rows
}

func (p *Presenter) Appointments(items []usecase.AppointmentResponseDTO) []AppointmentRow {
	rows := make([]AppointmentRow, 0, len(items))
	for _, item := range items {
		rows = append(rows, AppointmentRow{
			ID:        item.ID.String(),
			PatientID: item.PatientID.String(),
			DoctorID:  item.DoctorID.String(),
			Status:    string(item.Status),
			Date:      item.Date.Format(dateTimeLayout),
			Badge:     appointmentBadge(item.Status),
		})
	}
	return rows
}

func (p *Presenter) ScheduleSlots(items []usecase.ScheduleSlotResponseDTO) []ScheduleSlotRow {
	rows := make([]ScheduleSlotRow, 0, len(items))
	for _, item := range items {
		rows = append(rows, ScheduleSlotRow{
			ID:        item.ID.String(),
			DoctorID:  item.DoctorID.String(),
			Date:      item.Date.Format(dateLayout),
			StartTime: item.StartTime.Format(timeLayout),
			EndTime:   item.EndTime.Format(timeLayout),
			Cabinet:   item.Cabinet,
			IsBooked:  item.IsBooked,
			Badge:     bookedBadge(item.IsBooked),
		})
	}
	return rows
}

func (p *Presenter) Visits(items []usecase.VisitResponseDTO) []VisitRow {
	rows := make([]VisitRow, 0, len(items))
	for _, item := range items {
		rows = append(rows, VisitRow{
			ID:              item.ID.String(),
			AppointmentID:   item.AppointmentID.String(),
			Complaints:      item.Complaints,
			Diagnosis:       item.Diagnosis,
			Recommendations: item.Recommendations,
			Status:          string(item.Status),
			Badge:           visitBadge(item.Status),
		})
	}
	return rows
}

func appointmentBadge(status domain.AppointmentStatus) string {
	switch status {
	case domain.StatusCancelled:
		return "danger"
	case domain.StatusCompleted:
		return "muted"
	default:
		return "accent"
	}
}

func visitBadge(status domain.VisitStatus) string {
	if status == domain.VisitStatusCompleted {
		return "muted"
	}
	return "accent"
}

func bookedBadge(booked bool) string {
	if booked {
		return "muted"
	}
	return "accent"
}

func defaultDate() string {
	return time.Now().AddDate(0, 0, 1).Format(dateLayout)
}

func defaultDateTime() string {
	return time.Now().AddDate(0, 0, 1).Format(formTimeLayout)
}
