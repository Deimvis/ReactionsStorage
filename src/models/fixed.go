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
