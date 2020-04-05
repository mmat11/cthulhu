package bot

type Token string

type Config struct {
	Bot struct {
		AccessControl struct {
			Admin struct {
				IDs         []int    `yaml:"ids"`
				Permissions []string `yaml:"permissions"`
			}
			Moderator struct {
				IDs         []int    `yaml:"ids"`
				Permissions []string `yaml:"permissions"`
			}
		} `yaml:"access_control"`
	}
}

func (c *Config) CheckAdminPermissions(userID int, operation string) bool {
	var (
		userIDFound    = false
		operationFound = false
	)

	for _, id := range c.Bot.AccessControl.Admin.IDs {
		if id == userID {
			userIDFound = true
		}
	}
	for _, perm := range c.Bot.AccessControl.Admin.Permissions {
		if perm == operation {
			operationFound = true
		}
	}
	return userIDFound && operationFound
}
