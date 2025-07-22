package config

type Config struct {
	DbURL			string `json:"db_url"`
	CurrentUserName	string `json:"current_user_name"`
}

func (c Config) SetUser() error {
	c.CurrentUserName = username

	err := write(c)
	if err != nil {
		return err
	}

	return nil
}