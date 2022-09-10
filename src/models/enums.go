package models

type VisitStatus string

const (
	Created   VisitStatus = "Created"
	InProcess             = "InProcess"
	Done                  = "Done"
	Canceled              = "Canceled"
)

type ShiftStatus string

const (
	Planned  ShiftStatus = "Planned"
	Work     ShiftStatus = "Work"
	Finished ShiftStatus = "Finished"
)

func (s *ShiftStatus) String() string {
	return string(*s)
}
