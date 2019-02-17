package model_app

type Module struct {
	Id   int64 `json:"id,omitempty" bson:"_id,omitempty"`
	Name string
}
