CREATE TABLE public.users (
    id bigserial PRIMARY KEY,
    email character varying(255) UNIQUE NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    password bytea NOT NULL,
    activated bool NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	version integer NOT NULL DEFAULT 1
);

ALTER TABLE public.users OWNER TO auth_role;
