DELETE FROM "user_reaction"
WHERE "namespace_id" = $1 AND "entity_id" = $2 AND "reaction_id" = $3 AND "user_id" = $4