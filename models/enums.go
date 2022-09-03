package models

type VisitStatus string

const (
	Created   VisitStatus = "Created"
	InProcess             = "InProcess"
	Done                  = "Done"
	Canceled              = "Canceled"
)

type ServiceType string

const (
	Hair      ServiceType = "Hair"
	Beard                 = "Beard"
	HairBeard             = "HairBeard"
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
