package fsm

type State string

var (
	ShiftEnter State = "shiftEnter"
)

func (s *State) String() string {
	return string(*s)
}
