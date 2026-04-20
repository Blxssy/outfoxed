package domain

type DiceFace string

const (
	FaceFootprint string = "footprint"
	FaceEye       string = "eye"
)

type RollResult struct {
	Goal     GoalType `json:"goal"`
	Attempts int      `json:"attempts"`
	Faces    []string `json:"faces"`
	Success  bool     `json:"success"`
}

func RollForGoal(goal GoalType, rng RNG) RollResult {
	want := faceForGoal(goal)

	kept := make([]string, 0, 3)
	attempts := 0

	for attempts < 3 && len(kept) < 3 {
		attempts++

		need := 3 - len(kept)
		for i := 0; i < need; i++ {
			f := rollOneDice(rng)
			if f == want {
				kept = append(kept, f)
			}
		}
	}

	faces := make([]string, 0, 3)
	for i := 0; i < len(kept); i++ {
		faces = append(faces, want)
	}
	for len(faces) < 3 {
		faces = append(faces, oppositeFace(want))
	}

	return RollResult{
		Goal:     goal,
		Attempts: attempts,
		Faces:    faces,
		Success:  len(kept) == 3,
	}
}

func faceForGoal(goal GoalType) string {
	if goal == GoalClue {
		return FaceFootprint
	}
	return FaceEye
}

func rollOneDice(rng RNG) string {
	if rng.Intn(2) == 0 {
		return FaceFootprint
	}
	return FaceEye
}

func oppositeFace(f string) string {
	if f == FaceFootprint {
		return FaceEye
	}
	return FaceFootprint
}
