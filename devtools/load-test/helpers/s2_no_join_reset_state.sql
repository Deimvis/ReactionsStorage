DELETE FROM "user_reactions"
WHERE length(entity_id) != 36;
INSERT INTO "user_reactions"
SELECT
    namespace_id,
    entity_id,
    user_id,
    array_agg(reaction_id) AS reaction_ids
FROM "_sim_setup_200k_normdist_ur"
GROUP BY (namespace_id, entity_id, user_id)
;

DELETE FROM "reactions_count"
WHERE length(entity_id) != 36;
WITH tmp AS (
    SELECT
        namespace_id,
        entity_id,
        reaction_id,
        COUNT(*) AS reactions_cnt
    FROM "_sim_setup_200k_normdist_ur"
    GROUP BY namespace_id, entity_id, reaction_id
)
INSERT INTO "reactions_count"
SELECT
    namespace_id,
    entity_id,
    jsonb_object_agg(reaction_id, reactions_cnt) AS reactions_count
FROM tmp
GROUP BY namespace_id, entity_id
;