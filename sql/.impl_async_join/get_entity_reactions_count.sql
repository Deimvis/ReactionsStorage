SELECT
    "reactions_count"
FROM "entity_reactions"
WHERE "namespace_id__entity_id" = $1
