DELETE FROM "user_reaction"
WHERE "timestamp" > 1714124646;
REFRESH MATERIALIZED VIEW "entity_reactions";