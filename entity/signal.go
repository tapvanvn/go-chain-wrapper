package entity

type Signal struct {
	ItemName string            `json:"item_name" bson:"item_name"`
	Params   map[string]string `json:"signal" bson:"signal"`
}
