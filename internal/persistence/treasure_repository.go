package persistence

type TreasureRepository interface {
	ListTreasure(params map[string]string) []Treasure
	GetTreasureByID(treasureID int64) Treasure
}
type Treasure struct {
	TreasureID  int64
	Type        string
	Description string
}
