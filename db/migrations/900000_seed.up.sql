BEGIN;

INSERT INTO managers (id, email, password) 
VALUES ('4e0cac8a-eb8c-4c6c-8b07-2eaf831f4468', 'a@a.a', '$2a$10$blk0sKGaU8pNVuTXWqHrJOTW81/VaPDMNUJsww1Ai3ZpzPBcZbUoO');
INSERT INTO managers (id, email, password) 
VALUES ('614d4eb6-28e2-4771-ba61-24fa6182795b', 'b@a.a', '$2a$10$Kcyv01TybozWZ.REJKmt0.i3AGnl2jLnLKiqu4iwNhcf45rFpG32G');

INSERT INTO departments (id, name, manager_id)
VALUES ('e818019e-c8e0-4d20-8428-ddeeeafa3260', 'dept-a-1', '4e0cac8a-eb8c-4c6c-8b07-2eaf831f4468');
INSERT INTO departments (id, name, manager_id)
VALUES ('f6fd290b-d95c-42c3-99b4-c069f3a0049d', 'dept-b-1', '614d4eb6-28e2-4771-ba61-24fa6182795b');

COMMIT;

-- {
--     "email": "b@a.a",
--     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtYW5hZ2VyX2VtYWlsIjoiYkBhLmEiLCJtYW5hZ2VyX2lkIjoiNjE0ZDRlYjYtMjhlMi00NzcxLWJhNjEtMjRmYTYxODI3OTViIn0.NaIq0ELN_Yrv22CZiIWvFKuOQl0U1vfi-6rBHNauwSw"
-- }