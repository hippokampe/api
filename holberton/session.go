package holberton

func (hbtn *Holberton) newSession() (*holbertonSession, error) {
	ctx, err := hbtn.browser.NewContext()
	if err != nil {
		return nil, err
	}

	return &holbertonSession{
		User:           nil,
		BrowserContext: &ctx,
	}, nil
}

func (hbtn *Holberton) IsLogged(email string) bool {
	_, err := hbtn.getSession(email)
	return err == nil
}

func (hbtn *Holberton) getSession(email string) (*holbertonSession, error) {
	ctx, ok := hbtn.sessions[email]
	if !ok {
		return nil, ErrSessionNotExists
	}

	return ctx, nil
}

func (hbtn *Holberton) addSession(email string, ctx *holbertonSession) {
	hbtn.sessions[email] = ctx
}

func (hbtn *Holberton) deleteSession(email string) {
	delete(hbtn.sessions, email)
}
