WITH reaction_set_id AS (
    SELECT "reaction_set_id"
    FROM "namespace"
    WHERE "id" = $1
    LIMIT 1
),
reaction_ids AS (
    SELECT unnest("reaction_ids") as "reaction_id"
    FROM "reaction_set"
    WHERE "id" = (SELECT "reaction_set_id" FROM reaction_set_id)
)
SELECT
    "id",
    "short_name",
    "type",
    "code",
    "url"
FROM "reaction" as reaction
INNER JOIN "reaction_ids" as reaction_ids
    ON (reaction.id = reaction_ids.reaction_id)
