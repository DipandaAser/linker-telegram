package bot

type CustomReceiver struct {
	id string
}

func (c CustomReceiver) Recipient() string {
	return c.id
}
