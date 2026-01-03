package subjectsettings

func (ss *SubjectSettings) Modify(handler func(state *SubjectSettings) error) error {
	err := handler(ss)
	if err != nil {
		return err
	}
	ss.dirty = true
	return nil
}
