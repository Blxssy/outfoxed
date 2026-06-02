package domain

import (
	"fmt"
	"time"
)

type EventType string

const (
	EvGoalChosen       EventType = "goal_chosen"
	EvRolled           EventType = "rolled"
	EvRollFinished     EventType = "roll_finished"
	EvFoxMoved         EventType = "fox_moved"
	EvTurnEnded        EventType = "turn_ended"
	EvClueTaken        EventType = "clue_taken"
	EvAccused          EventType = "accused"
	EvGameFinished     EventType = "game_finished"
	EvPawnMoved        EventType = "pawn_moved"
	EvSuspectsRevealed EventType = "suspect_revealed"
)

const MaxJournalEntries = 100

type Event struct {
	Type EventType      `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}

type JournalEntry struct {
	ID        string         `json:"id"`
	Turn      int            `json:"turn"`
	Version   int            `json:"version"`
	Type      EventType      `json:"type"`
	Message   string         `json:"message"`
	Data      map[string]any `json:"data,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

func AppendEventsToJournal(st *GameState, events []Event, actor PlayerID) {
	if st == nil || len(events) == 0 {
		return
	}

	now := time.Now().UTC()

	for i, ev := range events {
		if !IsJournalEvent(ev.Type) {
			continue
		}

		message := FormatEventMessage(*st, ev, actor)
		if message == "" {
			continue
		}

		entry := JournalEntry{
			ID:        buildJournalEntryID(*st, ev, i),
			Turn:      st.Turn,
			Version:   st.Version,
			Type:      ev.Type,
			Message:   message,
			Data:      ev.Data,
			CreatedAt: now,
		}

		st.Journal = append(st.Journal, entry)
	}

	if len(st.Journal) > MaxJournalEntries {
		st.Journal = st.Journal[len(st.Journal)-MaxJournalEntries:]
	}
}

func buildJournalEntryID(st GameState, ev Event, index int) string {
	return fmt.Sprintf("%s:%d:%d:%d:%s", st.ID, st.Turn, st.Version, len(st.Journal)+index+1, ev.Type)
}

func FormatEventMessage(st GameState, ev Event, actor PlayerID) string {
	playerName := playerNameByID(st, actor)

	switch ev.Type {
	case EvGoalChosen:
		goal := eventString(ev, "goal")
		if goal == string(GoalClue) {
			return withPlayer(playerName, "выбрал цель: искать улику.")
		}
		if goal == string(GoalSuspect) {
			return withPlayer(playerName, "выбрал цель: проверить подозреваемых.")
		}
		return withPlayer(playerName, "выбрал цель хода.")

	case EvRolled:
		faces := eventStringSlice(ev, "faces")
		if len(faces) == 0 {
			return withPlayer(playerName, "бросил кубики.")
		}

		rollsUsed := eventInt(ev, "rollsUsed", 0)
		maxRolls := eventInt(ev, "maxRolls", 0)

		diceText := formatDiceFaces(faces)

		if rollsUsed > 0 && maxRolls > 0 {
			return withPlayer(
				playerName,
				fmt.Sprintf("бросил кубики %d/%d: %s.", rollsUsed, maxRolls, diceText),
			)
		}

		return withPlayer(
			playerName,
			fmt.Sprintf("бросил кубики: %s.", diceText),
		)

	case EvPawnMoved:
		return ""

	case EvRollFinished:
		return ""

	case EvClueTaken:
		trait := eventString(ev, "trait")
		result := eventString(ev, "result")
		return fmt.Sprintf("Найдена улика: %s — %s.", traitLabel(trait), traitResultLabel(result))

	case EvSuspectsRevealed:
		ids := eventStringSlice(ev, "ids")
		names := suspectNamesByIDs(st, ids)
		if len(names) > 0 {
			return "Открыты подозреваемые: " + joinHuman(names) + "."
		}
		return "Открыты новые подозреваемые."

	case EvFoxMoved:
		track := eventAny(ev, "track")
		return fmt.Sprintf("Лис продвинулся по следу. Текущая позиция: %v.", track)

	case EvTurnEnded:
		seat := eventInt(ev, "activeSeat", st.ActiveSeat)
		name := playerNameBySeat(st, seat)

		if name != "" {
			return fmt.Sprintf("Ход завершён. Расследование продолжает %s.", name)
		}

		return "Ход завершён. Очередь переходит следующему игроку."

	case EvAccused:
		suspectID := eventString(ev, "suspectId")
		correct := eventBool(ev, "correct")
		name := suspectNameByID(st, suspectID)

		if correct {
			return fmt.Sprintf("Сделано верное обвинение: %s оказался Лисом!", name)
		}
		return fmt.Sprintf("Сделано неверное обвинение: %s не был Лисом.", name)

	case EvGameFinished:
		result := eventString(ev, "result")
		if result == string(ResultWin) {
			return "Игра завершена победой команды."
		}
		if result == string(ResultLose) {
			return "Игра завершена поражением команды."
		}
		return "Игра завершена."

	case "turn_timed_out":
		seat := eventInt(ev, "seat", st.ActiveSeat)
		name := playerNameBySeat(st, seat)

		if name != "" {
			return fmt.Sprintf("%s не успел сходить. Ход доигрывает бот.", name)
		}

		return "Игрок не успел сходить. Ход доигрывает бот."

	case "clue_already_taken":
		return "Эта улика уже была найдена ранее."

	case "game_started":
		return "Игра началась. Расследование открыто!"

	default:
		return fmt.Sprintf("Событие: %s.", ev.Type)
	}
}

func IsJournalEvent(t EventType) bool {
	switch t {
	case
		EvGoalChosen,
		EvRolled,
		EvClueTaken,
		EvSuspectsRevealed,
		EvFoxMoved,
		EvTurnEnded,
		EvAccused,
		EvGameFinished:
		return true

	case "game_started", "turn_timed_out", "clue_already_taken":
		return true

	default:
		return false
	}
}

func formatDiceFaces(faces []string) string {
	if len(faces) == 0 {
		return ""
	}

	result := ""

	for i, face := range faces {
		if i > 0 {
			result += ", "
		}

		result += diceFaceLabel(face)
	}

	return result
}

func diceFaceLabel(face string) string {
	switch face {
	case string(FaceEye):
		return "👁️ Глаз"
	case string(FaceFootprint):
		return "👣 След"
	default:
		if face == "" {
			return "❔ Неизвестно"
		}
		return face
	}
}

func withPlayer(playerName string, message string) string {
	if playerName == "" {
		return "Игрок " + message
	}
	return playerName + " " + message
}

func playerNameByID(st GameState, userID PlayerID) string {
	if userID == "" {
		return ""
	}

	for _, p := range st.Players {
		if p.UserID == userID {
			return p.Name
		}
	}

	return ""
}

func suspectNameByID(st GameState, id string) string {
	for _, suspect := range st.Suspects {
		if suspect.ID == id {
			if suspect.Code != "" {
				return string(suspect.Code)
			}
			return suspect.ID
		}
	}

	return id
}

func suspectNamesByIDs(st GameState, ids []string) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, suspectNameByID(st, id))
	}
	return out
}

func joinHuman(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}

	result := ""
	for i, item := range items {
		switch {
		case i == 0:
			result = item
		case i == len(items)-1:
			result += " и " + item
		default:
			result += ", " + item
		}
	}

	return result
}

func traitLabel(value string) string {
	switch value {
	case string(ClueTraitGlasses):
		return "Очки"
	case string(ClueTraitHat):
		return "Шляпа"
	case string(ClueTraitScarf):
		return "Шарф"
	case string(ClueTraitUmbrella):
		return "Зонтик"
	case string(ClueTraitBag):
		return "Сумка"
	case string(ClueTraitBoots):
		return "Ботинки"
	default:
		if value == "" {
			return "Неизвестно"
		}
		return value
	}
}

func traitResultLabel(value string) string {
	switch value {
	case string(TraitYes):
		return "есть"
	case string(TraitNo):
		return "нет"
	default:
		if value == "" {
			return "неизвестно"
		}
		return value
	}
}

func eventAny(ev Event, key string) any {
	if ev.Data == nil {
		return nil
	}
	return ev.Data[key]
}

func eventString(ev Event, key string) string {
	if ev.Data == nil {
		return ""
	}

	v, ok := ev.Data[key]
	if !ok || v == nil {
		return ""
	}

	return fmt.Sprint(v)
}

func eventBool(ev Event, key string) bool {
	if ev.Data == nil {
		return false
	}

	v, ok := ev.Data[key]
	if !ok {
		return false
	}

	switch x := v.(type) {
	case bool:
		return x
	case string:
		return x == "true" || x == "yes"
	default:
		return false
	}
}

func eventStringSlice(ev Event, key string) []string {
	if ev.Data == nil {
		return nil
	}

	v, ok := ev.Data[key]
	if !ok || v == nil {
		return nil
	}

	switch x := v.(type) {
	case []string:
		return x

	case []interface{}:
		out := make([]string, 0, len(x))
		for _, item := range x {
			out = append(out, fmt.Sprint(item))
		}
		return out

	default:
		return nil
	}
}

func playerNameBySeat(st GameState, seat int) string {
	for _, p := range st.Players {
		if p.Seat == seat {
			return p.Name
		}
	}

	return ""
}

func eventInt(ev Event, key string, fallback int) int {
	if ev.Data == nil {
		return fallback
	}

	v, ok := ev.Data[key]
	if !ok || v == nil {
		return fallback
	}

	switch x := v.(type) {
	case int:
		return x
	case int8:
		return int(x)
	case int16:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float64:
		return int(x)
	case float32:
		return int(x)
	default:
		return fallback
	}
}
