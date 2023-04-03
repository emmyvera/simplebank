CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "acct_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_acct_id" bigint NOT NULL,
  "to_acct_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("acct_id");

CREATE INDEX ON "transfers" ("from_acct_id");

CREATE INDEX ON "transfers" ("to_acct_id");

CREATE INDEX ON "transfers" ("from_acct_id", "to_acct_id");

COMMENT ON COLUMN "entries"."amount" IS 'It can be negetive or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'It must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("acct_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_acct_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_acct_id") REFERENCES "accounts" ("id");
