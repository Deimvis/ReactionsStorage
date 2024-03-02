SELECT
    "reaction_id",
    COUNT(*) as "count"
FROM "user_reaction"
WHERE "namespace_id" = $1 AND "entity_id" = $2
GROUP BY "reaction_id"
