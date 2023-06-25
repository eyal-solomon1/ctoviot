CREATE TABLE "users" (
  "username" varchar UNIQUE PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "balance" bigint NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "videos" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "video_name" varchar NOT NULL,
  "video_identifier" varchar  NOT NULL,
  "video_length" bigint NOT NULL,
  "video_remote_path" varchar UNIQUE NOT NULL,
  "video_decs" varchar NOT NULL,
  "created_at" timestamptz  NOT NULL DEFAULT (now()),
  CONSTRAINT "unique_owner_video" UNIQUE ("owner", "video_identifier", "video_name")
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "video_name" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "videos" ("owner");

CREATE INDEX ON "entries" ("username");

ALTER TABLE "videos" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
