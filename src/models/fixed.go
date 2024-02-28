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
