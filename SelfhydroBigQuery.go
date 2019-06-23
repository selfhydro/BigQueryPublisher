package bigQueryPublisher

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

type PubSubMessage struct {
	Data        []byte            `json:"data"`
	Attributes  map[string]string `json:"attributes"`
	PublishTime string            `json:"publishTime"`
}

type SeflhydroState struct {
	AmbientTemperature          float64 `json:"ambientTemperature"`
	AmbientHumidity             float64 `json:"ambientHumidity"`
	WaterTemperature            float64 `json:"waterTemperature"`
	WaterElectricalConductivity float64 `json:"waterElectricalConductivity"`
	Time                        string  `json:"time"`
	deviceId                    string
}

type SelfhydroStateTable struct {
	WaterLevel                  bigquery.NullFloat64 `bigquery:"waterLevel"`
	AmbientTemperature          bigquery.NullFloat64 `bigquery:"ambientTemperature"`
	AmbientHumidity             bigquery.NullFloat64 `bigquery:"ambientHumidity"`
	WaterTemperature            bigquery.NullFloat64 `bigquery:"waterTemperature"`
	WaterElectricalConductivity bigquery.NullFloat64 `bigquery:"waterElectricalConductivity"`
	Time                        civil.DateTime       `bigquery:"time"`
	DeviceId                    string               `bigquery:"deviceId"`
}

var bigqueryClient *bigquery.Client
var ctx context.Context
var projectId = os.Getenv("GCP_PROJECT")
var datasetId = os.Getenv("DATASET_ID")
var tableId = os.Getenv("TABLE_ID")

func init() {
	var err error
	ctx = context.Background()
	bigqueryClient, err = bigquery.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("couldnt init big query client: %v", err.Error())
	}
}

func TransferStateToBigQuery(ctx context.Context, m PubSubMessage) error {
	state := DeseraliseState(m.Data)
	state.deviceId = m.Attributes["deviceId"]
	log.Printf("received state from: %s, the current ambient temperature is: %f as of %v", state.deviceId, state.AmbientTemperature, state.Time)
	saveToStateTable(state)
	return nil
}

func DeseraliseState(data []byte) SeflhydroState {
	var state = SeflhydroState{}
	err := json.Unmarshal(data, &state)
	if err != nil {
		log.Fatalf("can't decode state from message: %v", err.Error())
	}
	return state
}

func saveToStateTable(selfhydroState SeflhydroState) error {
	tableInserter := bigqueryClient.Dataset(datasetId).Table(tableId).Inserter()
	time, err := time.Parse("20060102150405", selfhydroState.Time)
	if err != nil {
		log.Printf("couldn't parse state time: %v", err.Error())
	}
	log.Printf("time parsed: %v", time)
	datetime := civil.DateTimeOf(time)
	log.Printf("time converted to civil: %v", datetime)

	temperature := convertFloatToBigQueryFloat(selfhydroState.AmbientTemperature)
	humidity := convertFloatToBigQueryFloat(selfhydroState.AmbientHumidity)

	waterTemperature := convertFloatToBigQueryFloat(selfhydroState.WaterTemperature)
	waterElectricalConductivity := convertFloatToBigQueryFloat(selfhydroState.WaterElectricalConductivity)
	state := []*SelfhydroStateTable{
		{DeviceId: selfhydroState.deviceId, Time: datetime, AmbientTemperature: temperature, AmbientHumidity: humidity, WaterTemperature: waterTemperature, WaterElectricalConductivity: waterElectricalConductivity},
	}
	if err := tableInserter.Put(ctx, state); err != nil {
		log.Printf("cant insert state into big query: ")
		log.Printf("%v", err.Error())
		return err
	}
	return nil
}

func convertFloatToBigQueryFloat(value float64) bigquery.NullFloat64 {
	if value == 0 {
		return bigquery.NullFloat64{}
	}
	return bigquery.NullFloat64{Float64: value, Valid: true}
}
