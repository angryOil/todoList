SET statement_timeout = 0;

--bun:split

CREATE TABLE "public"."todo"
(
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(50) UNIQUE NOT NULL,
    content         VARCHAR(300),
    order_num       int,
    is_deleted      bool,
    created_at      timestamptz,
    last_updated_at timestamptz
);
