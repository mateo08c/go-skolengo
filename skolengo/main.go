package skolengo

func NewClient(seleniumURL string) (*Client, error) {
	return &Client{
		SeleniumURL: seleniumURL,

		AutoLogin: true,
	}, nil
}
