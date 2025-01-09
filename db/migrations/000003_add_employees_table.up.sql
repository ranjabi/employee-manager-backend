BEGIN;

CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE employees (
    id uuid             PRIMARY KEY DEFAULT gen_random_uuid(),
    identity_number     text NOT NULL UNIQUE,
    name text NOT NULL,
    employee_image_uri  text NOT NULL,
    gender              gender NOT NULL,
    department_id       uuid NOT NULL,

    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT
);

COMMIT;