CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "customer_name" varchar NOT NULL,
  "ordered_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "items" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "quantity" int NOT NULL,
  "order_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

ALTER TABLE "items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;

CREATE INDEX ON "items" ("order_id");

CREATE INDEX ON "items" ("description");