package skolengo

func NewClient(username string, password string) (*Client, error) {
	return &Client{
		Username:  username,
		Password:  password,
		AutoLogin: true,
	}, nil
}
