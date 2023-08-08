package entity

type Tasks struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Title    string `json:"title" bson:"title"`
	ActiveAt string `json:"activeAt" bson:"activeAt"`
	Status   string `json:"status" bson:"status"`
}
