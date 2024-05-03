DELETE FROM "user_reaction"
WHERE "timestamp" > 1714124646;
VACUUM "user_reaction";