package types

import "time"

type WarningType string

const (
	SpecialWeatherStatement WarningType = "Special Weather Statement"
	FlashFloodWarning       WarningType = "Flash Flood Warning"
	SevereThunderstorm      WarningType = "Severe Thunderstorm Warning"
	Tornado                 WarningType = "Tornado Warning"
)

type NotificationType string

const (
	New    NotificationType = "New"
	Update NotificationType = "Update"
	Cancel NotificationType = "Cancel"
)

type WarningMap map[string]*Warning

type Warning struct {
	Type          WarningType
	Notification  NotificationType
	References    []Reference
	OldWarning    *Warning
	Timer         *time.Timer `json:"-"`
	MessageID     string
	ID            string
	NWSCenter     string
	AffectedAreas []string
	Expires       time.Time
	TornadoSource string
	WindThreat    string
	HailThreat    string
	Path          [][]float64
}

type Reference struct {
	ID     string    `json:"@id"`
	ID0    string    `json:"identifier"`
	Sender string    `json:"sender"`
	Sent   time.Time `json:"sent"`
}

type ActiveResponse struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Features []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Geometry struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			ID       string `json:"@id"`
			Type     string `json:"@type"`
			ID0      string `json:"id"`
			AreaDesc string `json:"areaDesc"`
			Geocode  struct {
				Same []string `json:"SAME"`
				Ugc  []string `json:"UGC"`
			} `json:"geocode"`
			AffectedZones []string    `json:"affectedZones"`
			References    []Reference `json:"references"`
			Sent          time.Time   `json:"sent"`
			Effective     time.Time   `json:"effective"`
			Onset         time.Time   `json:"onset"`
			Expires       time.Time   `json:"expires"`
			Ends          time.Time   `json:"ends"`
			Status        string      `json:"status"`
			MessageType   string      `json:"messageType"`
			Category      string      `json:"category"`
			Severity      string      `json:"severity"`
			Certainty     string      `json:"certainty"`
			Urgency       string      `json:"urgency"`
			Event         string      `json:"event"`
			Sender        string      `json:"sender"`
			SenderName    string      `json:"senderName"`
			Headline      string      `json:"headline"`
			Description   string      `json:"description"`
			Instruction   string      `json:"instruction"`
			Response      string      `json:"response"`
			Parameters    struct {
				AWIPSidentifier        []string    `json:"AWIPSidentifier"`
				WMOidentifier          []string    `json:"WMOidentifier"`
				EventMotionDescription []string    `json:"eventMotionDescription"`
				WindThreat             []string    `json:"windThreat"`
				MaxWindGust            []string    `json:"maxWindGust"`
				HailThreat             []string    `json:"hailThreat"`
				MaxHailSize            []string    `json:"maxHailSize"`
				TornadoDetection       []string    `json:"tornadoDetection"`
				Blockchannel           []string    `json:"BLOCKCHANNEL"`
				EASORG                 []string    `json:"EAS-ORG"`
				Vtec                   []string    `json:"VTEC"`
				EventEndingTime        []time.Time `json:"eventEndingTime"`
			} `json:"parameters"`
		} `json:"properties"`
	} `json:"features"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
}
