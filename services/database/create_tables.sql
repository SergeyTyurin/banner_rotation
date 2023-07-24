CREATE SEQUENCE Banner_id_seq AS integer;
CREATE TABLE IF NOT EXISTS "Banners"(
    "id" integer NOT NULL DEFAULT nextval('Banner_id_seq'),
    "info" text,
    PRIMARY KEY ("id")
);
ALTER SEQUENCE Banner_id_seq OWNED BY "Banners"."id";

CREATE SEQUENCE Slot_id_seq AS integer;
CREATE TABLE IF NOT EXISTS "Slots"(
    "id" integer NOT NULL DEFAULT nextval('Slot_id_seq'),
    "info" text,
    PRIMARY KEY ("id")
);
ALTER SEQUENCE Slot_id_seq OWNED BY "Slots"."id";

CREATE SEQUENCE Group_id_seq AS integer;
CREATE TABLE IF NOT EXISTS "Groups"(
    "id" integer NOT NULL DEFAULT nextval('Group_id_seq'),
    "info" text,
    PRIMARY KEY ("id")
);
ALTER SEQUENCE Group_id_seq OWNED BY "Groups"."id";

CREATE TABLE IF NOT EXISTS "Statistic"(
    "slot_id" integer,
    "group_id" integer,
    "banner_id" integer,
    "display_count" integer NOT NULL DEFAULT 0,
    "click_count" integer NOT NULL DEFAULT 0,
    FOREIGN KEY ("slot_id") REFERENCES "Slots" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("group_id") REFERENCES "Groups" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("banner_id") REFERENCES "Banners" ("id") ON DELETE CASCADE 
);