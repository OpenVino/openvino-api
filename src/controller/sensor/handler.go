package sensor

import (
	customHTTP "github.com/openvino/openvino-api/src/http"
	"github.com/openvino/openvino-api/src/model"
	"github.com/openvino/openvino-api/src/repository"
	"log"
	"net/http"
)

type QueryData struct {
	Harvest   string  `json:"year"`
	Month     string  `json:"month"`
	Day		  string  `json:"day"`
}

func GetSensorRecords(w http.ResponseWriter, r *http.Request) {

	var query string = "min(timestamp) as timestamp, sensor_id,"+
		"avg(humidity2) as humidity2, avg(humidity1) as humidity1,"+
		"avg(humidity05) as humidity05, avg(humidity005) as humidity005,"+
		"max(wind_velocity) as wind_velocity, max(wind_gust) as wind_gust,"+
		"avg(wind_direction) as wind_direction, avg(pressure) as pressure,"+
		"max(rain) as rain, avg(temperature) as temperature,"+
		"avg(humidity) as humidity, max(irradiance_ir) as irradiance_ir,"+
		"max(irradiance_uv) as irradiance_uv, max(irradiance_vi) as irradiance_vi"

	var params = QueryData{}
	params.Harvest = r.URL.Query().Get("year")
	params.Month = r.URL.Query().Get("month")
	params.Day = r.URL.Query().Get("day")

	log.Println(params)

	records := []model.SensorRecord{};

	if params.Day == "" && params.Month == "" && params.Harvest != "" {
		repository.DB.Select(query).
			Where("YEAR(timestamp) = ?", params.Harvest).
			Group("DAY(timestamp), MONTH(timestamp), sensor_id").
			Find(&records);
	} else if params.Day == "" && params.Month != "" && params.Harvest != "" {
		repository.DB.Select(query).
			Where("MONTH(timestamp) = ? AND YEAR(timestamp) = ?", params.Month, params.Harvest).
			Group("DAY(timestamp), sensor_id").
			Find(&records);
	} else if params.Day != "" && params.Month != "" && params.Harvest != "" {
		repository.DB.
			Where("DAY(timestamp) = ? AND MONTH(timestamp) = ? AND YEAR(timestamp) = ?", params.Day, params.Month, params.Harvest).
			Find(&records);
	} else {
		sensordataCs := model.SensorRecord{}
		sensordataPv := model.SensorRecord{}
		sensordataMo := model.SensorRecord{}
		sensordataMe := model.SensorRecord{}
		repository.DB.Where("sensor_id = ?", "petit-verdot").Order("timestamp desc").Limit(1).Find(&sensordataPv)
		repository.DB.Where("sensor_id = ?", "cabernet-sauvignon").Order("timestamp desc").Limit(1).Find(&sensordataCs)
		repository.DB.Where("sensor_id = ?", "malbec-este").Order("timestamp desc").Limit(1).Find(&sensordataMe)
		repository.DB.Where("sensor_id = ?", "malbec-oeste").Order("timestamp desc").Limit(1).Find(&sensordataMo)
		records = []model.SensorRecord{sensordataCs, sensordataPv, sensordataMo, sensordataMe}
	}
	customHTTP.ResponseJSON(w, records)
	return
}

func GetSensorHashes(w http.ResponseWriter, r *http.Request) {

	var params = QueryData{}
	params.Harvest = r.URL.Query().Get("year")
	params.Month = r.URL.Query().Get("month")
	params.Day = r.URL.Query().Get("day")

	var hashes []string

	if params.Day == "" && params.Month == "" && params.Harvest != "" {
		repository.DB.Table("sensor_records").
			Where("YEAR(timestamp) = ?", params.Harvest).Order("timestamp desc").
			Pluck("hash", &hashes);

	} else if params.Day == "" && params.Month != "" && params.Harvest != "" {
		repository.DB.Table("sensor_records").
			Where("MONTH(timestamp) = ? AND YEAR(timestamp) = ?", params.Month, params.Harvest).
			Order("timestamp desc").
			Pluck("hash", &hashes);
	} else if params.Day != "" && params.Month != "" && params.Harvest != "" {
		repository.DB.Table("sensor_records").
			Where("DAY(timestamp) = ? AND MONTH(timestamp) = ? AND YEAR(timestamp) = ?", params.Day, params.Month, params.Harvest).
			Order("timestamp desc").
			Pluck("hash", &hashes);
	} else {
		var sensordataCs string
		var sensordataPv string
		var sensordataMo string
		var sensordataMe string
		repository.DB.Select("hash").Where("sensor_id = ?", "petit-verdot").Order("timestamp desc").Limit(1).First(&sensordataPv)
		repository.DB.Where("sensor_id = ?", "cabernet-sauvignon").Order("timestamp desc").Limit(1).First(&sensordataCs)
		repository.DB.Where("sensor_id = ?", "malbec-este").Order("timestamp desc").Limit(1).First(&sensordataMe)
		repository.DB.Where("sensor_id = ?", "malbec-oeste").Order("timestamp desc").Limit(1).First(&sensordataMo)
		hashes = []string{sensordataCs, sensordataPv, sensordataMo, sensordataMe}
	}
	customHTTP.ResponseJSON(w, hashes)
	return
}
