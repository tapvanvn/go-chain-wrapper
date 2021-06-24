package entity

type Param struct {
	Type  string `json:"type" bson:"type"`
	Value []byte `json:"value" bson:"value"`
}
