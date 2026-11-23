-- Table: public.last_email

-- DROP TABLE IF EXISTS public.last_email;

CREATE TABLE IF NOT EXISTS public.last_email
(
    message_id character varying(200) COLLATE pg_catalog."default" NOT NULL,
    internal_date bigint,
    owner_id integer NOT NULL, -- multi-owner: one checkpoint row per owner
    CONSTRAINT last_email_owner_id_key UNIQUE (owner_id),
    CONSTRAINT last_email_owner_id_fkey FOREIGN KEY (owner_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.last_email
    OWNER to postgres;

-- Table: public.processed_gmail_messages

-- DROP TABLE IF EXISTS public.processed_gmail_messages;

CREATE TABLE IF NOT EXISTS public.processed_gmail_messages
(
    id bigint NOT NULL DEFAULT nextval('processed_gmail_messages_id_seq'::regclass),
    gmail_message_id character varying(200) COLLATE pg_catalog."default" NOT NULL,
    internal_date bigint NOT NULL,
    transaction_id bigint,
    status character varying(20) COLLATE pg_catalog."default" NOT NULL,
    error_message text COLLATE pg_catalog."default",
    processed_at timestamp without time zone NOT NULL DEFAULT now(),
    owner_id integer NOT NULL, -- multi-owner
    CONSTRAINT processed_gmail_messages_pkey PRIMARY KEY (id),
    -- multi-owner: Gmail message IDs are unique per mailbox, not globally
    -- (replaces uq_gmail_message)
    CONSTRAINT uq_gmail_message_owner UNIQUE (owner_id, gmail_message_id),
    CONSTRAINT processed_gmail_owner_id_fkey FOREIGN KEY (owner_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.processed_gmail_messages
    OWNER to postgres;
-- Index: idx_processed_gmail_internal_date

-- DROP INDEX IF EXISTS public.idx_processed_gmail_internal_date;

CREATE INDEX IF NOT EXISTS idx_processed_gmail_internal_date
    ON public.processed_gmail_messages USING btree
    (internal_date ASC NULLS LAST)
    WITH (fillfactor=100, deduplicate_items=True)
    TABLESPACE pg_default;
-- Index: idx_processed_gmail_status

-- DROP INDEX IF EXISTS public.idx_processed_gmail_status;

CREATE INDEX IF NOT EXISTS idx_processed_gmail_status
    ON public.processed_gmail_messages USING btree
    (status COLLATE pg_catalog."default" ASC NULLS LAST)
    WITH (fillfactor=100, deduplicate_items=True)
    TABLESPACE pg_default;

-- Table: public.users_token  (unchanged; public.users now lives in initial_table.sql)

-- DROP TABLE IF EXISTS public.users_token;

CREATE TABLE IF NOT EXISTS public.users_token
(
    user_id integer NOT NULL DEFAULT nextval('users_token_user_id_seq'::regclass),
    refresh_token character varying(50) COLLATE pg_catalog."default",
    token character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT users_token_pkey PRIMARY KEY (user_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users_token
    OWNER to postgres;