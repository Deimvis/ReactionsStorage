package configs

type Simulation struct {
	Seed   int64
	Rules  Rules  `yaml:"rules" json:"rules"`
	Server Server `yaml:"server" json:"server"`
}

type Rules struct {
	TurnCount int `yaml:"turn_count" json:"turn_count"`
	TurnDurS  int `yaml:"turn_duration_s" json:"turn_duration_s"`
	UserCount int `yaml:"user_count" json:"user_count"`
	Users     []struct {
		Count  int `yaml:"count" json:"count"`
		Screen struct {
			VisibleEntitiesCount int `yaml:"visible_entities_count" json:"visible_entities_count"`
		} `yaml:"screen" json:"screen"`
	} `yaml:"screen" json:"screen"`
	Topics []struct {
		Count       int    `yaml:"count" json:"count"`
		NamespaceId string `yaml:"namespace_id" json:"namespace_id"`
		Size        int    `yaml:"size" json:"size"`
	} `yaml:"topics" json:"topics"`
}

type Server struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
	SSL  bool   `yaml:"ssl" json:"ssl"`
}
