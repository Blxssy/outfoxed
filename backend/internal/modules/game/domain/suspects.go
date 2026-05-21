package domain

import "strconv"

type SuspectCard struct {
	ID       string        `json:"id"`
	Code     SuspectCode   `json:"code"`
	Revealed bool          `json:"revealed"`
	Excluded bool          `json:"excluded"`
	Traits   SuspectTraits `json:"traits"`
}

type SuspectTraits struct {
	Glasses  TraitValue `json:"glasses"`
	Hat      TraitValue `json:"hat"`
	Scarf    TraitValue `json:"scarf"`
	Umbrella TraitValue `json:"umbrella"`
	Bag      TraitValue `json:"bag"`
	Boots    TraitValue `json:"boots"`
}

type SuspectCardView struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Revealed bool           `json:"revealed"`
	Excluded bool           `json:"excluded"`
	Traits   *SuspectTraits `json:"traits,omitempty"`
}

type ClueToken struct {
	ID        string      `json:"id"`
	Trait     ClueTrait   `json:"trait"`
	Revealed  bool        `json:"revealed"`
	Result    *TraitValue `json:"result,omitempty"`
	BoardCell int         `json:"boardCell"`
}

type ClueTokenView struct {
	ID        string      `json:"id"`
	Revealed  bool        `json:"revealed"`
	Trait     *ClueTrait  `json:"trait,omitempty"`
	Result    *TraitValue `json:"result,omitempty"`
	BoardCell int         `json:"boardCell"`
}

func suspectIDFromIndex(i int) string {
	return "suspect_" + strconv.Itoa(i+1)
}
