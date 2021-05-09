package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/parkedwards/go-trainer-api/pkg/models"
)

// This client only reads from one "table" - appointments.json
const connectionString = "./pkg/database/appointments.json"

// This "DBClient" module acts as a database client,
// abstracting the underlying datastore implementation (in this case, a JSON file :)
type DBClient struct{}

func New() *DBClient {
	return &DBClient{}
}

// Not the most efficient, as every query will parse the JSON + load into memory
// but the goal was not to implement an efficient datastore here
func (dbc *DBClient) OpenDatabaseConnection(connectionString string) []models.Appointment {
	raw, err := ioutil.ReadFile(connectionString)
	if err != nil {
		fmt.Println(err)
	}

	data := []models.Appointment{}
	err = json.Unmarshal([]byte(raw), &data)
	if err != nil {
		fmt.Println(err)
	}

	return data
}

func (dbc *DBClient) Select(q *models.Query) []models.Appointment {
	data := dbc.OpenDatabaseConnection(connectionString)

	results := []models.Appointment{}

	for _, a := range data {
		if a.TrainerId == q.TrainerId {
			results = append(results, a)
		}
	}
	return results
}

func (dbc *DBClient) Insert(appointment *models.Appointment) {
	data := dbc.OpenDatabaseConnection(connectionString)
	appointmentId := int64(uuid.New().ID())
	appointment.Id = &appointmentId
	data = append(data, *appointment)

	raw, _ := json.MarshalIndent(data, "", "	")
	ioutil.WriteFile(connectionString, raw, 0644)
}
