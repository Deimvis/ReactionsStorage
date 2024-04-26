-- $1 - namespace_id
-- $2 - entity_id
-- $3 - reaction_id
-- $4 - user_id
-- $5 - last_update_ts (unix timestapm in seconds)

INSERT INTO "user_reactions" ("namespace_id", "entity_id", "user_id", "reaction_ids")
VALUES
($1, $2, $4, ARRAY[$3::TEXT])
ON CONFLICT ("namespace_id", "entity_id", "user_id") DO UPDATE
SET "reaction_ids" = array_append("user_reactions"."reaction_ids", $3::TEXT),
    "last_update_ts" = $5

;

INSERT INTO "reactions_count" ("namespace_id", "entity_id", "reactions_count") VALUES
($1, $2, jsonb_build_object($3::TEXT, 1))
ON CONFLICT ("namespace_id", "entity_id") DO UPDATE
SET "reactions_count" = jsonb_set("reactions_count"."reactions_count", ARRAY[$3::TEXT], to_jsonb(COALESCE(("reactions_count"."reactions_count" ->> $3::TEXT)::int, 0) + 1))
