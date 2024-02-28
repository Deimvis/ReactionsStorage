package sql

var AddUserReaction = ReadQueryFile("add_user_reaction.sql")
var GetMaxUniqueReactions = ReadQueryFile("get_max_unique_reactions.sql")
var GetMutuallyExclusiveReactions = ReadQueryFile("get_mutually_exclusive_reactions.sql")
var GetUniqueEntityReactions = ReadQueryFile("get_unique_entity_reactions.sql")
var GetUniqueEntityUserReactions = ReadQueryFile("get_unique_entity_user_reactions.sql")
var GetEntityReactionsCount = ReadQueryFile("get_entity_reactions_count.sql")
var InitReactionsStorage = ReadQueryFile("init_reactions_storage.sql")
