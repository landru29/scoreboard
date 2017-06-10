package chronometer

import (
	"time"
	"encoding/json"

	"github.com/landru29/scoreboard/src/server/ws"
)

// Create a new chronometer
func Create(client *ws.Client) *Chronometer {
	c := &Chronometer {
		Running: false,
		Client: client,
		Tick : time.NewTicker(time.Millisecond * 500),
	}

	defer func() {
		c.Tick.Stop()
		client.Connexion.Close()
	}()

	go func() {
		for {
			select {
			case control := <-c.Control:
				now := time.Now()
				status = Status {
					Ellapsed: c.Ellapsed.String(),
					Now: now.Format("2006-01-02 15:04:05"),
					Status: control,
					UUID: c.Client.UUID,
				}
				if control == "stop" {
					if c.Running {
						projection := now.Add(c.Ellapsed)
						c.Ellapsed += projection.Sub(c.Starting)
						c.Running = false
						if data, err := json.Marshal(status); err != nil {
							c.Client.SendMessage([]byte("{\"status\":\"fatal\"}"))
						} else {
							c.Client.SendMessage(data)
						}
					}
					
				}

				if control == "start" {
					c.Starting = time.Now()
					c.Running = true
					if data, err := json.Marshal(status); err != nil {
						c.Client.SendMessage([]byte("{\"status\":\"fatal\"}"))
					} else {
						c.Client.SendMessage(data)
					}
				}

				if control == "close" {
					if data, err := json.Marshal(status); err != nil {
						c.Client.SendMessage([]byte("{\"status\":\"fatal\"}"))
					} else {
						c.Client.SendMessage(data)
					}
					return
				}

			case <-c.Tick.C:
				if (!c.Running) {
					status.Status = "idle"
				}
				if data, err := json.Marshal(status); err != nil {
					c.Client.SendMessage([]byte("{\"status\":\"fatal\"}"))
				} else {
					c.Client.SendMessage(data)
				}
		}

	return c
}
