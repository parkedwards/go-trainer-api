package models

type Query struct {
	TrainerId int64 `json:"trainer_id,omitempty"`
}

type Appointment struct {
	Id        int64  `json:"id,omitempty"`
	TrainerId int64  `json:"trainer_id"`
	UserId    int64  `json:"user_id"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
}

type Availability struct {
	TrainerId int64  `json:"trainer_id"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
}
