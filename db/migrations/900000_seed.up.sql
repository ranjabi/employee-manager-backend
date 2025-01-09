BEGIN;

INSERT INTO managers (id, email, password) 
VALUES ('4e0cac8a-eb8c-4c6c-8b07-2eaf831f4468', 'a@a.a', '$2a$10$blk0sKGaU8pNVuTXWqHrJOTW81/VaPDMNUJsww1Ai3ZpzPBcZbUoO');
INSERT INTO managers (id, email, password) 
VALUES ('614d4eb6-28e2-4771-ba61-24fa6182795b', 'b@a.a', '$2a$10$Kcyv01TybozWZ.REJKmt0.i3AGnl2jLnLKiqu4iwNhcf45rFpG32G');

INSERT INTO departments (id, name, manager_id)
VALUES ('e818019e-c8e0-4d20-8428-ddeeeafa3260', 'dept-a-1', '4e0cac8a-eb8c-4c6c-8b07-2eaf831f4468');
INSERT INTO departments (id, name, manager_id)
VALUES ('f6fd290b-d95c-42c3-99b4-c069f3a0049d', 'dept-b-1', '614d4eb6-28e2-4771-ba61-24fa6182795b');

INSERT INTO employees (
    id
    ,identity_number
    ,name
    ,employee_image_uri
    ,gender
    ,department_id
) 
VALUES (
    '626f360a-1b72-4f1c-9c37-5bc6d8adc1d1'
    ,'iden1emp'
    ,'emp1name'
    ,'http://www.google.com'
    ,'male'
    ,'e818019e-c8e0-4d20-8428-ddeeeafa3260'
);
INSERT INTO employees (
    id
    ,identity_number
    ,name
    ,employee_image_uri
    ,gender
    ,department_id
) 
VALUES (
    'eeb4eae1-d3a4-41ec-80a5-bdec6b5b6b0b'
    ,'iden1emp2'
    ,'emp1name2'
    ,'http://www.google.com'
    ,'male'
    ,'e818019e-c8e0-4d20-8428-ddeeeafa3260'
);
INSERT INTO employees (
    id
    ,identity_number
    ,name
    ,employee_image_uri
    ,gender
    ,department_id
) 
VALUES (
    'b48cfe2b-c225-4ad2-b1c3-8acf3249fac6'
    ,'iden2emp'
    ,'emp2name'
    ,'http://www.google.com'
    ,'female'
    ,'f6fd290b-d95c-42c3-99b4-c069f3a0049d'
);

COMMIT;

-- {
--     "email": "b@a.a",
--     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtYW5hZ2VyX2VtYWlsIjoiYkBhLmEiLCJtYW5hZ2VyX2lkIjoiNjE0ZDRlYjYtMjhlMi00NzcxLWJhNjEtMjRmYTYxODI3OTViIn0.NaIq0ELN_Yrv22CZiIWvFKuOQl0U1vfi-6rBHNauwSw"
-- }