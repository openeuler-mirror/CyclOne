package mysql

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"idcos.io/cloudboot/model"
)

func loadCriticalSensorDatas() (items []*model.SensorData) {
	data, _ := ioutil.ReadFile("./testdata/ipmi_critical.json")
	_ = json.Unmarshal(data, &items)
	return items
}

func loadNominalSensorDatas() (items []*model.SensorData) {
	data, _ := ioutil.ReadFile("./testdata/ipmi_nominal.json")
	_ = json.Unmarshal(data, &items)
	return items
}

func loadWarningSensorDatas() (items []*model.SensorData) {
	data, _ := ioutil.ReadFile("./testdata/ipmi_warning.json")
	_ = json.Unmarshal(data, &items)
	return items
}
func Test_healthStatus(t *testing.T) {
	Convey("ipmi检查设备健康状况", t, func() {

		criticals := loadCriticalSensorDatas()
		So(criticals, ShouldNotBeEmpty)
		So(new(InspectionJob).healthStatus(criticals), ShouldEqual, model.HealthStatusCritical)

		warnings := loadWarningSensorDatas()
		So(warnings, ShouldNotBeEmpty)
		So(new(InspectionJob).healthStatus(warnings), ShouldEqual, model.HealthStatusWarning)

		nominals := loadNominalSensorDatas()
		So(nominals, ShouldNotBeEmpty)
		So(new(InspectionJob).healthStatus(nominals), ShouldEqual, model.HealthStatusNominal)
	})
}
