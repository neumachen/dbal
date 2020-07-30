-- Migration goes here.
CREATE TABLE customers(
  customer_id SERIAL primary key,
  first_name  VARCHAR not null,
  last_name   VARCHAR not null,
  address     JSONB   not null,
  created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
  updated_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL
);
