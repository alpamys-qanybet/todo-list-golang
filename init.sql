-- status
CREATE TABLE public.status (
    name character varying(120) NOT NULL
);

ALTER TABLE public.status OWNER TO postgres;

INSERT INTO status VALUES ('created'), ('in_progress'), ('paused'), ('done'), ('deleted');

-- task
CREATE TABLE public.task (
    id integer NOT NULL,
    name character varying(255),
    description character varying(1200),
    status character varying(120)
);

ALTER TABLE public.task OWNER TO postgres;

ALTER TABLE public.task ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.task_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 65535
    CACHE 1
);

ALTER TABLE ONLY public.status ADD CONSTRAINT status_key PRIMARY KEY (name);

CREATE INDEX fki_status_fk ON public.task USING btree (status);

ALTER TABLE ONLY public.task ADD CONSTRAINT status_fk FOREIGN KEY (status) REFERENCES public.status(name) NOT VALID;
