package entity

type Signal struct {
	ItemName string             `json:"item_name" bson:"item_name"`
	Params   map[string][]Param `json:"signal" bson:"signal"`
}
