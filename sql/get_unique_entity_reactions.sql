SELECT
    "reaction_ids"
FROM "user_reactions"
WHERE "namespace_id" = $1 AND "entity_id" = $2
