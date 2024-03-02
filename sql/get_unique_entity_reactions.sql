SELECT
    DISTINCT "reaction_id"
FROM "user_reaction"
WHERE "namespace_id" = $1 AND "entity_id" = $2
