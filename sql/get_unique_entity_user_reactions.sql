SELECT
    DISTINCT "reaction_id"
FROM "user_reaction"
WHERE "namespace_id" = $1 AND "entity_id" = $2 AND "user_id" = $3
