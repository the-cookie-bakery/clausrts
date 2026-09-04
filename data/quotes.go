// Package data holds static content: the original German quotes from Vom Kriege
// and their English glosses, surfaced by the teaching layer.
package data

// Quote pairs a Clausewitz concept with the original German text and an English
// gloss for the reader.
type Quote struct {
	Concept string
	German  string
	English string
}

// Quotes is the curated set of teaching quotes, keyed by concept.
var Quotes = map[string]Quote{
	"trinity": {
		Concept: "Trinity",
		German:  "eine wunderliche Dreifaltigkeit, zusammengesetzt aus der ursprünglichen Gewaltsamkeit seines Elementes, dem Haß und der Feindschaft ... aus dem Spiel der Wahrscheinlichkeiten und des Zufalls ... und aus der untergeordneten Natur eines politischen Werkzeuges",
		English: "A wonderful trinity, composed of the original violence of its element (hatred and enmity), the play of probabilities and chance, and its subordination as an instrument of policy.",
	},
	"policy": {
		Concept: "War as policy",
		German:  "Der Krieg ist eine bloße Fortsetzung der Politik mit anderen Mitteln.",
		English: "War is merely a continuation of policy by other means.",
	},
	"friction": {
		Concept: "Friction",
		German:  "Die Gefahr im Kriege gehört zur Friktion desselben ...",
		English: "Danger in war belongs to its friction: real war is not war on paper.",
	},
	"schwerpunkt": {
		Concept: "Center of gravity",
		German:  "daß man also die strategischen Wirkungen auf gewisse Schwerpunkte zurückführen kann",
		English: "Strategic effects can be traced back to certain centers of gravity.",
	},
	"culmination": {
		Concept: "Culminating point",
		German:  "Über den Kulminationspunkt des Sieges ...",
		English: "Beyond the culminating point of victory, the offensive exhausts itself.",
	},
}
