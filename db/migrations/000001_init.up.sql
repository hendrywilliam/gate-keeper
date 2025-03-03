CREATE TABLE "concert" (
    "id" serial PRIMARY KEY,
    "name" varchar(255) UNIQUE,
    "artist_id" integer,
    "venue_id" integer,
    "date" int,
    "limit" int,
    "created_at" int,
    "updated_at" int
);

CREATE TABLE "ticket" (
    "id" serial PRIMARY KEY,
    "serial_number" uuid DEFAULT gen_random_uuid(),
    "concert_id" integer,
    "ticket_category_id" integer,
    "created_at" int,
    "updated_at" int
);

CREATE TABLE "ticket_category" (
    "id" serial PRIMARY KEY,
    "concert_id" integer,
    "description" varchar(255),
    "price" decimal(10,2),
    "start_date" int,
    "end_date" int,
    "created_at" int,
    "updated_at" int
);

ALTER TABLE "ticket_category" ADD FOREIGN KEY ("concert_id") REFERENCES "concert" ("id");

ALTER TABLE "ticket" ADD FOREIGN KEY ("concert_id") REFERENCES "concert" ("id");

ALTER TABLE "ticket" ADD FOREIGN KEY ("ticket_category_id") REFERENCES "ticket_category" ("id");
