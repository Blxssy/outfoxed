package domain

const (
	FaceFootprint = "footprint"
	FaceEye       = "eye"
	DiceCount     = 3
	MaxRolls      = 3
)

type RollResult struct {
	Goal     GoalType `json:"goal"`
	Attempts int      `json:"attempts"`
	Faces    []string `json:"faces"`
	Success  bool     `json:"success"`
}

func faceForGoal(goal GoalType) string {
	if goal == GoalClue {
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

func rollOneFace(rng RNG) string {
	if rng.Intn(2) == 0 {
		return FaceFootprint
	}
	return FaceEye
}

func rollInitialFaces(rng RNG) []string {
	return []string{
		rollOneFace(rng),
		rollOneFace(rng),
		rollOneFace(rng),
	}
}

func rerollFaces(current []string, keep []bool, rng RNG) []string {
	out := make([]string, len(current))
	copy(out, current)

	for i := range out {
		if i >= len(keep) || !keep[i] {
			out[i] = rollOneFace(rng)
		}
	}

	return out
}

func isRollSuccessful(goal GoalType, faces []string) bool {
	want := faceForGoal(goal)
	for _, f := range faces {
		if f != want {
			return false
		}
	}
	return len(faces) == DiceCount
}

// Старый авто-режим оставляем для timeout/fallback.
func RollForGoal(goal GoalType, rng RNG) RollResult {
	faces := rollInitialFaces(rng)
	attempts := 1

	success := isRollSuccessful(goal, faces)
	for attempts < MaxRolls && !success {
		keep := make([]bool, len(faces))
		want := faceForGoal(goal)

		for i, f := range faces {
			if f == want {
				keep[i] = true
			}
		}

		faces = rerollFaces(faces, keep, rng)
		attempts++
		success = isRollSuccessful(goal, faces)
	}

	return RollResult{
		Goal:     goal,
		Attempts: attempts,
		Faces:    faces,
		Success:  success,
	}
}
