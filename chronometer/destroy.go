package chronometer

func (c *Chronometer) Destroy() {
	c.Control <- "destroy"
}
