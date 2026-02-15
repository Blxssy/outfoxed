package domain

type SuspectState struct {
	ID       int  `json:"id"`
	Revealed bool `json:"revealed"`
	Excluded bool `json:"excluded"`
}

func NewSuspects(total int) []SuspectState {
	res := make([]SuspectState, total)
	for i := 0; i < total; i++ {
		res[i] = SuspectState{
			ID:       i,
			Revealed: false,
			Excluded: false,
		}
	}
	return res
}
