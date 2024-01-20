package models

type Config struct {
	DB struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Name string `yaml:"name"`
	}
	SERVER struct {
		Host string `yaml:"host"`
	}
	DEBUG struct {
		Switch string `yaml:"switch"`
		Login  string `yaml:"login"`
		Cors   string `yaml:"cors"`
	}
	Bot struct {
		APIUrl    string `yaml:"apiUrl"`
		Token     string `yaml:"token"`
		GetUpdate string `yaml:"getUpdate"`
		Methods   struct {
			SendMessage string `yaml:"sendMessage"`
		}
	}
}
