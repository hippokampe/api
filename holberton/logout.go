package holberton

func (h *Holberton) logout() error {
	var err error

	_, err = h.page.Goto(BaseUrl + "/users/my_profile")
	if err != nil {
		return err
	}

	err = h.page.Click(".text-right > a:nth-child(1)")
	if err != nil {
		return err
	}

	h.InternalStatus.Logged = false

	return nil
}
