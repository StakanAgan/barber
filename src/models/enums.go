package models

type VisitStatus string

const (
	Created       VisitStatus = "Created"
	Accepted      VisitStatus = "Accepted"
	InProcess     VisitStatus = "InProcess"
	Done          VisitStatus = "Done"
	VisitCanceled VisitStatus = "Canceled"
)

type ShiftStatus string

const (
	Planned       ShiftStatus = "Planned"
	Work          ShiftStatus = "Work"
	Finished      ShiftStatus = "Finished"
	ShiftCanceled ShiftStatus = "Canceled"
)

func (s *VisitStatus) String() string {
	return string(*s)
}

func (s *ShiftStatus) String() string {
	return string(*s)
}
