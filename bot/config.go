package bot

type Token string

type TaskArg struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type TaskArgs []struct {
	Arg TaskArg `yaml:"arg"`
}

type Config struct {
	Bot struct {
		Token         Token
		AccessControl struct {
			Mods struct {
				IDs []int `yaml:"ids"`
			} `yaml:"mods"`
			Groups []struct {
				Group struct {
					ID             int64    `yaml:"id"`
					URL            string   `yaml:"url"`
					WelcomeMessage string   `yaml:"welcome_message"`
					CrossPostTags  []string `yaml:"crosspost_tags"`
					Admin          struct {
						IDs         []int    `yaml:"ids"`
						Permissions []string `yaml:"permissions"`
					} `yaml:"admin"`
				} `yaml:"group"`
			} `yaml:"groups"`
		} `yaml:"access_control"`
		Tasks []struct {
			Task struct {
				Name string   `yaml:"name"`
				Cron string   `yaml:"cron"`
				Args TaskArgs `yaml:"args"`
			} `yaml:"task"`
		} `yaml:"tasks"`
	} `yaml:"bot"`
}

func (c *Config) hasPermissions(chatID int64, userID int, operation string) bool {
	var (
		userIDFound    = false
		operationFound = false
	)

	for _, g := range c.Bot.AccessControl.Groups {
		if g.Group.ID == chatID {
			for _, id := range g.Group.Admin.IDs {
				if id == userID {
					userIDFound = true
				}
			}
			for _, perm := range g.Group.Admin.Permissions {
				if perm == operation {
					operationFound = true
				}
			}
		}
	}
	return userIDFound && operationFound
}

func (c *Config) isMod(userID int) bool {
	for _, id := range c.Bot.AccessControl.Mods.IDs {
		if id == userID {
			return true
		}
	}
	return false
}
