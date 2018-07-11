-- +goose Up
CREATE TABLE IF NOT EXISTS article (
  id integer NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone,
  user_id integer,
  title text,
  body text,
  tags text
);

CREATE SEQUENCE article_id_seq
  START with 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

ALTER SEQUENCE article_id_seq OWNED BY article.id;

ALTER TABLE ONLY article ALTER COLUMN id SET DEFAULT nextval('article_id_seq'::regclass);

ALTER TABLE ONLY article ADD CONSTRAINT article_pkey PRIMARY KEY (id);

CREATE INDEX idx_article_deleted_at ON article USING btree (deleted_at);

-- +goose Down
DROP TABLE article;