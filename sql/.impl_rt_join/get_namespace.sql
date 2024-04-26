SELECT
    "id",
    "reaction_set_id",
    "max_uniq_reactions",
    "mutually_exclusive_reactions"
FROM "namespace"
WHERE "id" = $1
