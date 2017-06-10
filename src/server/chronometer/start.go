package chronometer

// Start the chronometer
func (c *Chronometer) Start() {
	c.Control <- "start"
}
