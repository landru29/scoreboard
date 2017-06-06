package chronometer

// Pause the chronometer
func (c *Chronometer) Stop() {
	c.Control <- "stop"
}
