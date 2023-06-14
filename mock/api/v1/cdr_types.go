package v1

type DataStatus string

const (
	DataStatusUnknown    DataStatus = "00"
	DataStatusReceived   DataStatus = "01"
	DataStatusProcessing DataStatus = "02"
	DataStatusProcessed  DataStatus = "03"
)

type Data struct {
	UID    string     `json:"uniqueid"`
	Status DataStatus `json:"Status"`
}

type DataSet struct {
	Results []Data `json:"results"`
}

type Request struct {
	Action  string  `json:"action"`
	DataSet DataSet `json:"CdrdataSet"`
}
