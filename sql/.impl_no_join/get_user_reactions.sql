SELECT
    "namespace_id", 
    "entity_id",
    unnest("reaction_ids") as "reaction_id",
    "user_id"
FROM "user_reactions"
