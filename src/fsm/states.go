package fsm

type State string

var (
	ShiftEnter           State = "shiftEnter"
	ServiceEnterTitle    State = "serviceEnterTitle"
	ServiceEnterDuration State = "serviceEnterDuration"
	ServiceEnterPrice    State = "serviceEnterPrice"
)

func (s *State) String() string {
	return string(*s)
}
