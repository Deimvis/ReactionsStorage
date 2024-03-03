package models

func (resp *ReactionsGETResponse200) Code() int {
	return 200
}

func (resp *ReactionsPOSTResponse200) Code() int {
	return 200
}

func (resp *ReactionsPOSTResponse403) Code() int {
	return 403
}

func (resp *ReactionsDELETEResponse200) Code() int {
	return 200
}

func (resp *ReactionsDELETEResponse403) Code() int {
	return 403
}

func (resp *ConfigurationPOSTResponse200) Code() int {
	return 200
}

func (resp *ConfigurationPOSTResponse422) Code() int {
	return 422
}

func (resp *AvailableReactionsGETResponse200) Code() int {
	return 200
}

func (resp *AvailableReactionsGETResponse404) Code() int {
	return 404
}
