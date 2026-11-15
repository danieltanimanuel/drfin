-- Table: public.transaction_types

-- DROP TABLE IF EXISTS public.transaction_types;

CREATE TABLE IF NOT EXISTS public.transaction_types
(
    id integer NOT NULL DEFAULT nextval('transaction_types_id_seq'::regclass),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    CONSTRAINT transaction_types_pkey PRIMARY KEY (id),
    CONSTRAINT transaction_types_name_key UNIQUE (name)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transaction_types
    OWNER to postgres;

-- Table: public.categories

-- DROP TABLE IF EXISTS public.categories;

CREATE TABLE IF NOT EXISTS public.categories
(
    id integer NOT NULL DEFAULT nextval('categories_id_seq'::regclass),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (id),
    CONSTRAINT categories_name_key UNIQUE (name)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.categories
    OWNER to postgres;

-- Table: public.account_groups

-- DROP TABLE IF EXISTS public.account_groups;

CREATE TABLE IF NOT EXISTS public.account_groups
(
    id integer NOT NULL DEFAULT nextval('account_groups_id_seq'::regclass),
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    CONSTRAINT account_groups_pkey PRIMARY KEY (id),
    CONSTRAINT account_groups_name_key UNIQUE (name)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.account_groups
    OWNER to postgres;

-- Table: public.transactions

-- DROP TABLE IF EXISTS public.transactions;

CREATE TABLE IF NOT EXISTS public.transactions
(
    id integer NOT NULL DEFAULT nextval('transactions_id_seq'::regclass),
    parent_id integer NOT NULL,
    trade_date date NOT NULL,
    type_id integer NOT NULL,
    category_id integer NOT NULL,
    account_id integer NOT NULL,
    account_dest_id integer,
    amount numeric(24,2) NOT NULL,
    description character varying(1000) COLLATE pg_catalog."default",
    status character varying(15) COLLATE pg_catalog."default" NOT NULL,
    source character varying(15) COLLATE pg_catalog."default" NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    last_updated_at timestamp without time zone NOT NULL,
    confirmed_at timestamp without time zone,
    CONSTRAINT transactions_pkey PRIMARY KEY (id),
    CONSTRAINT transactions_account_dest_id_fkey FOREIGN KEY (account_dest_id)
        REFERENCES public.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transactions_account_id_fkey FOREIGN KEY (account_id)
        REFERENCES public.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transactions_category_id_fkey FOREIGN KEY (category_id)
        REFERENCES public.categories (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT transactions_type_id_fkey FOREIGN KEY (type_id)
        REFERENCES public.transaction_types (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transactions
    OWNER to postgres;

-- Table: public.accounts

-- DROP TABLE IF EXISTS public.accounts;

CREATE TABLE IF NOT EXISTS public.accounts
(
    id integer NOT NULL DEFAULT nextval('accounts_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    balance numeric(24,2) NOT NULL,
    account_group_id integer,
    status character varying(1) COLLATE pg_catalog."default" NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    closed_at timestamp without time zone,
    account_number character varying(20) COLLATE pg_catalog."default",
    credit_limit numeric(24,2),
    CONSTRAINT accounts_pkey PRIMARY KEY (id),
    CONSTRAINT accounts_name_key UNIQUE (name),
    CONSTRAINT accounts_account_group_id_fkey FOREIGN KEY (account_group_id)
        REFERENCES public.account_groups (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.accounts
    OWNER to postgres;

-- Table: public.owned_account

-- DROP TABLE IF EXISTS public.owned_account;

CREATE TABLE IF NOT EXISTS public.owned_account
(
    account_label character varying(30) COLLATE pg_catalog."default" NOT NULL,
    account_number character varying(20) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT owned_account_pkey PRIMARY KEY (account_label)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.owned_account
    OWNER to postgres;