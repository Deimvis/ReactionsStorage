package configs

type Simulation struct {
	Seed                  int64
	Rules                 Rules                  `yaml:"rules" json:"rules"`
	Server                Server                 `yaml:"server" json:"server"`
	PrometheusPushgateway *PrometheusPushgateway `yaml:"prometheus_pushgateway,omitempty" json:"prometheus_pushgateway"`
}

type Rules struct {
	Turns struct {
		Count    int `yaml:"count" json:"count"`
		MinDurMs int `yaml:"min_duration_ms" json:"min_duration_ms"`
	} `yaml:"turns" json:"turns"`

	Users struct {
		Count           int     `yaml:"count" json:"count"`
		TurnStartSkewMs int     `yaml:"turn_start_skew_ms" json:"turn_start_skew_ms"`
		IdTemplate      *string `yaml:"id_template" json:"id_template"`
		Screen          struct {
			VisibleEntitiesCount int `yaml:"visible_entities_count" json:"visible_entities_count"`
		} `yaml:"screen" json:"screen"`
		App struct {
			Background struct {
				RefreshReactions struct {
					TimerInTurns int `yaml:"timer_in_turns" json:"timer_in_turns"`
				} `yaml:"refresh_reactions" json:"refresh_reactions"`
			} `yaml:"background" json:"background"`
		} `yaml:"app" json:"app"`
		ActionProbs ActionProbs `yaml:"action_probs" json:"action_probs"`
	} `yaml:"users" json:"users"`

	Topics []struct {
		Count          int    `yaml:"count" json:"count"`
		NamespaceId    string `yaml:"namespace_id" json:"namespace_id"`
		Size           int    `yaml:"size" json:"size"`
		ShufflePerUser bool   `yaml:"shuffle_per_user" json:"shuffle_per_user"`
	} `yaml:"topics" json:"topics"`
}

type Server struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
	SSL  bool   `yaml:"ssl" json:"ssl"`
}

type PrometheusPushgateway struct {
	Host          string `yaml:"host" json:"host"`
	Port          int    `yaml:"port" json:"port"`
	SSL           bool   `yaml:"ssl" json:"ssl"`
	PushIntervalS int    `yaml:"push_interval_s" json:"push_interval_s"`
}

type ActionProbs struct {
	SwitchTopic    uint `yaml:"switch_topic" json:"switch_topic"`
	Scroll         uint `yaml:"scroll" json:"scroll"`
	AddReaction    uint `yaml:"add_reaction" json:"add_reaction"`
	RemoveReaction uint `yaml:"remove_reaction" json:"remove_reaction"`
	Quit           uint `yaml:"quit" json:"quit"`
}
