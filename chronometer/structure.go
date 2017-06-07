package chronometer

import (
	"time"

	"github.com/landru29/scoreboard/ws"
)

//Chronometer define a chronometer
type Chronometer struct {
	Starting time.Time
	Running  bool
	Ellapsed time.Duration
	Tick     *time.Ticker
	Control  chan string
	Client   *ws.Client
}

type Status struct {
	Ellapsed string `json:"ellapsed"`
	Status   string `json:"status"`
	Now      string `json:"now"`
	UUID     string `json:"uuid"`
}
