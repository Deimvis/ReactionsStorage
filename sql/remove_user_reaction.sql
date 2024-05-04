-- $1 - namespace_id
-- $2 - entity_id
-- $3 - reaction_id
-- $4 - user_id

UPDATE "user_reactions"
SET "reaction_ids" = array_remove("reaction_ids", $3)
WHERE "namespace_id" = $1 AND "entity_id" = $2 AND "user_id" = $4

;

UPDATE "reactions_count"
SET "reactions_count" = jsonb_set("reactions_count", ARRAY[$3::TEXT], to_jsonb(("reactions_count" ->> $3::TEXT)::int - 1))
WHERE "namespace_id" = $1 AND "entity_id" = $2
