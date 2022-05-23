-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TODO: Add type foreign key on all logbookentries

CREATE TABLE IF NOT EXISTS users
(
    id               uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at       TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    display_name     varchar(256) not null,
    email            varchar(256) not null,
    user_external_id varchar(256),
    activated        bool         NOT NULL             default false,
    activationCode   varchar(256)
);

CREATE TABLE IF NOT EXISTS organizational_unit_types
(
    id   uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(256) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS specialties
(
    id   uuid                NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(256) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS organizational_units
(
    id           uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    display_name VARCHAR(256) NOT NULL             DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,
    parent_id    uuid references organizational_units (id),
    metadata     jsonb        NOT NULL             default '{}',
    type_id      uuid         NOT NULL references organizational_unit_types (id)
);

CREATE TABLE IF NOT EXISTS organizational_unit_specialties
(
    unit_id      uuid NOT NULL references organizational_units (id) ON DELETE CASCADE,
    specialty_id uuid NOT NULL references specialties (id) ON DELETE CASCADE,
    PRIMARY KEY (unit_id, specialty_id)
);

CREATE TABLE IF NOT EXISTS roles
(
    id           serial PRIMARY KEY,
    name         VARCHAR(256) NOT NULL UNIQUE,
    display_name VARCHAR(256) NOT NULL DEFAULT ''
);

CREATE TABLE user_organizational_unit_roles
(
    unit_id uuid REFERENCES organizational_units (id) ON DELETE CASCADE,
    user_id uuid REFERENCES users (id) ON DELETE CASCADE,
    role_id integer REFERENCES roles (id) ON DELETE CASCADE,
    PRIMARY KEY (unit_id, user_id, role_id)
);

CREATE INDEX unit_id_idx ON user_organizational_unit_roles (unit_id);
CREATE INDEX role_id_idx ON user_organizational_unit_roles (role_id);
CREATE INDEX user_id_idx ON user_organizational_unit_roles (user_id);

-- TODO: Add specialities and user_specialty
CREATE TABLE IF NOT EXISTS practical_activity_types
(
    id           serial PRIMARY KEY,
    display_name varchar(256) NOT NULL,
    name         varchar(256) NOT NULL UNIQUE, -- can be used to as filter when listing logbookentries.
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS specialties_activity_types
(
    activity_type_id integer NOT NULL references practical_activity_types (id) ON DELETE CASCADE,
    specialty_id     uuid    NOT NULL references specialties (id) ON DELETE CASCADE,
    PRIMARY KEY (activity_type_id, specialty_id)
);

CREATE TABLE iF NOT EXISTS procedures
(
    id                     uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    display_name           varchar(256) NOT NULL,
    created_at             TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    annotations            jsonb        NOT NULL,
    organizational_unit_id uuid         NOT NULL references organizational_units (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS examinations
(
    id                         uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    display_name               varchar(256) NOT NULL,
    annotations                jsonb        NOT NULL,
    created_at                 TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    department_id              uuid         NOT NULL references organizational_units (id) ON DELETE CASCADE,
    practical_activity_type_id integer      NOT NULL references practical_activity_types (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS examination_activities
(
    id                 uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_user_id   uuid         NOT NULL references users (id) ON DELETE CASCADE,
    supervisor_user_id uuid references users (id) ON DELETE CASCADE,
    examination_id     uuid         NOT NULL references examinations (id) ON DELETE CASCADE,
    display_name       varchar(256) NOT NULL             DEFAULT '',
    created_at         TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    annotations        jsonb        NOT NULL
);

CREATE TABLE IF NOT EXISTS examinations_activity_reviews
(
    id                 uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    supervisor_user_id uuid         NOT NULL references users (id) ON DELETE CASCADE,
    display_name       varchar(256) NOT NULL,
    created_at         TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    annotations        jsonb        NOT NULL,
    comment            VARCHAR(1024)
);

CREATE TABLE IF NOT EXISTS examinations_activities_reviews
(
    id                               uuid        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    activity_author_user_id          uuid        NOT NULL references users (id) ON DELETE CASCADE,
    activity_reviewer_user_id        uuid        NOT NULL references users (id) ON DELETE CASCADE,
    examination_activities_id        uuid        NOT NULL references examination_activities (id) ON DELETE CASCADE,
    examinations_activity_reviews_id uuid        references examinations_activity_reviews (id) ON DELETE SET NULL,
    created_at                       TIMESTAMPTZ NOT NULL             DEFAULT NOW(),
    resident_updated_at              TIMESTAMPTZ,
    supervisor_updated_at            TIMESTAMPTZ,
    is_reviewed                      boolean     NOT NULL             DEFAULT false
);

CREATE TABLE IF NOT EXISTS surgery_diagnosis
(
    id            uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    bodypart      varchar(256) NOT NULL,
    diagnose_name varchar(256) NOT NULL             DEFAULT '',
    diagnose_code varchar(256) NOT NULL             DEFAULT '',
    extra_code    varchar(256) NOT NULL             DEFAULT ''
);

CREATE TABLE IF NOT EXISTS surgery_methods
(
    id            uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    method_name   varchar(256) NOT NULL             DEFAULT '',
    method_code   varchar(256) NOT NULL             DEFAULT '',
    approach_name varchar(256) NOT NULL             DEFAULT ''
);

CREATE TABLE IF NOT EXISTS surgeries
(
    id                uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    diagnose_id       uuid NOT NULL references surgery_diagnosis (id),
    method_id         uuid NOT NULL references surgery_methods (id),
    surgery_specialty uuid NOT NULL REFERENCES specialties (id) ON DELETE CASCADE,
    unique (diagnose_id, method_id, surgery_specialty)
);

CREATE TABLE IF NOT EXISTS orthopedic_surgery_activities
(
    id                         uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurred_at                TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    created_at                 TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    case_notes                 VARCHAR(256) NOT NULL,
    patient_age                integer      NOT NULL,
    patient_gender             VARCHAR(256) NOT NULL,
    resident_id                uuid         NOT NULL references users (id),
    supervisor_id              uuid REFERENCES users (id),
    operator_id                uuid REFERENCES users (id),
    assistant_id               uuid REFERENCES users (id),
    comments                   VARCHAR(256) NOT NULL             DEFAULT '',
    complications              VARCHAR(256) NOT NULL             DEFAULT '',
    annotations                jsonb        NOT NULL             DEFAULT '{}'::jsonb,
    dops_requested             boolean      NOT NULL             default false,
    review_requested           boolean      NOT NULL             default false,
    has_dops_connection        boolean      NOT NULL             DEFAULT false,
    practical_activity_type_id integer      NOT NULL references practical_activity_types (id) ON DELETE CASCADE,
    in_progress                boolean      NOT NULL             default false,
    active_step                integer      NOT NULL             default 0,
    completed_step             integer      NOT NULL             default 0
);

CREATE TABLE IF NOT EXISTS orthopedic_surgery_activities_surgeries
(
    id                             uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    surgery_id                     uuid NOT NULL references surgeries (id),
    orthopedic_surgery_activity_id uuid NOT NULL references orthopedic_surgery_activities (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orthopedic_surgeries_activity_review
(
    id                             uuid         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    orthopedic_surgery_activity_id uuid         NOT NULL REFERENCES orthopedic_surgery_activities (id) UNIQUE,
    created_at                     TIMESTAMPTZ  NOT NULL             DEFAULT NOW(),
    updated_at                     TIMESTAMPTZ,
    signed_at                      TIMESTAMPTZ,
    occurred_at                    TIMESTAMPTZ  NOT NULL,
    case_notes                     VARCHAR(256) NOT NULL,
    patient_age                    integer      NOT NULL,
    patient_gender                 VARCHAR(256) NOT NULL,
    annotations                    jsonb        NOT NULL             DEFAULT '{}'::jsonb,
    resident_id                    uuid         NOT NULL references users (id),
    supervisor_id                  uuid         NOT NULL references users (id),
    operator_id                    uuid         NOT NULL references users (id),
    assistant_id                   uuid         NOT NULL references users (id),
    comments                       VARCHAR(256) NOT NULL             DEFAULT '',
    complications                  VARCHAR(256) NOT NULL             DEFAULT '',
    review_comment                 VARCHAR(256) NOT NULL             DEFAULT '',
    requested                      boolean      NOT NULL             default false,
    in_progress                    boolean      NOT NULL             default false,
    active_step                    integer      NOT NULL             default 0,
    completed_step                 integer      NOT NULL             default 0
);

CREATE TABLE IF NOT EXISTS orthopedic_surgeries_activity_review_surgeries
(
    id                                      uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    surgery_id                              uuid NOT NULL references surgeries (id),
    orthopedic_surgeries_activity_review_id uuid NOT NULL references orthopedic_surgeries_activity_review (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS evaluation_forms
(
    id            serial PRIMARY KEY,
    department_id uuid references organizational_units (id) NOT NULL,
    name          VARCHAR(256)                              NOT NULL,
    annotations   jsonb                                     NOT NULL DEFAULT '{}'::jsonb,
    difficulty    text[]                                    NOT NULL DEFAULT '{}',
    citations     text[]                                    NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS dops_evaluations
(
    id                             uuid         NOT NULL PRIMARY KEY                                             DEFAULT uuid_generate_v4(),
    orthopedic_surgery_activity_id uuid         references orthopedic_surgery_activities (id) ON DELETE SET NULL DEFAULT NULL UNIQUE,
    resident_id                    uuid         NOT NULL references users (id),
    supervisor_id                  uuid         NOT NULL references users (id),
    occurred_at                    TIMESTAMPTZ  NOT NULL,
    case_notes                     VARCHAR(256) NOT NULL,
    patient_age                    integer      NOT NULL,
    patient_gender                 VARCHAR(256) NOT NULL,
    difficulty                     VARCHAR(256) NOT NULL                                                         DEFAULT '',
    department_id                  uuid         references organizational_units (id) ON DELETE SET NULL          DEFAULT NULL,
    annotations                    jsonb        NOT NULL                                                         DEFAULT '{}'::jsonb,
    is_evaluated                   boolean      NOT NULL                                                         DEFAULT false,
    created_at                     TIMESTAMPTZ  NOT NULL                                                         DEFAULT NOW(),
    in_progress                    boolean      NOT NULL                                                         default false,
    active_step                    integer      NOT NULL                                                         default 0,
    completed_step                 integer      NOT NULL                                                         default 0
);

CREATE TABLE IF NOT EXISTS dops_evaluations_surgeries
(
    id                 uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    surgery_id         uuid NOT NULL references surgeries (id),
    dops_evaluation_id uuid NOT NULL references dops_evaluations (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mini_cex_evaluations
(
    id             uuid         NOT NULL PRIMARY KEY                                    DEFAULT uuid_generate_v4(),
    resident_id    uuid         NOT NULL references users (id),
    supervisor_id  uuid         NOT NULL references users (id),
    occurred_at    TIMESTAMPTZ  NOT NULL,
    difficulty     VARCHAR(256) NOT NULL                                                DEFAULT '',
    area           VARCHAR(256) NOT NULL                                                DEFAULT '',
    focuses        text[]       NOT NULL                                                DEFAULT '{}',
    department_id  uuid         references organizational_units (id) ON DELETE SET NULL DEFAULT NULL,
    annotations    jsonb        NOT NULL                                                DEFAULT '{}'::jsonb,
    is_evaluated   boolean      NOT NULL                                                DEFAULT false,
    created_at     TIMESTAMPTZ  NOT NULL                                                DEFAULT NOW(),
    in_progress    boolean      NOT NULL                                                default false,
    active_step    integer      NOT NULL                                                default 0,
    completed_step integer      NOT NULL                                                default 0
);

CREATE TABLE IF NOT EXISTS mini_cex_areas
(
    id            serial PRIMARY KEY,
    name          VARCHAR(256) NOT NULL,
    department_id uuid         NOT NULL references organizational_units (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mini_cex_focuses
(
    id   serial PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS notifications
(
    id          uuid        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    annotations jsonb       NOT NULL             DEFAULT '{}'::jsonb,
    user_id     uuid        NOT NULL references users (id),
    created_at  TIMESTAMPTZ NOT NULL             DEFAULT NOW(),
    seen_at     TIMESTAMPTZ                      DEFAULT null
);

CREATE TABLE IF NOT EXISTS logbook_entries
(
    id                    uuid        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurred_at           TIMESTAMPTZ NOT NULL,
    orthopedic_surgery_id uuid references orthopedic_surgery_activities (id) ON DELETE CASCADE UNIQUE,
    examination_id        uuid references examinations (id) ON DELETE CASCADE UNIQUE,
    procedure_id          uuid references procedures (id) ON DELETE CASCADE UNIQUE
);

CREATE TABLE IF NOT EXISTS assessments
(
    id                           uuid        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurred_at                  TIMESTAMPTZ NOT NULL,
    orthopedic_surgery_review_id uuid references orthopedic_surgeries_activity_review ON DELETE CASCADE UNIQUE,
    dops_id                      uuid references dops_evaluations ON DELETE CASCADE UNIQUE,
    mini_cex_id                  uuid references mini_cex_evaluations ON DELETE CASCADE UNIQUE
);

CREATE TABLE IF NOT EXISTS activities
(
    id               uuid        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    occurred_at      TIMESTAMPTZ NOT NULL,
    logbook_entry_id uuid references logbook_entries (id) ON DELETE CASCADE UNIQUE,
    assessment_id    uuid references assessments (id) ON DELETE CASCADE UNIQUE
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activities;
DROP TABLE IF EXISTS assessments;
DROP TABLE IF EXISTS logbook_entries;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS mini_cex_focuses;
DROP TABLE IF EXISTS mini_cex_areas;
DROP TABLE IF EXISTS mini_cex_evaluations;
DROP TABLE IF EXISTS dops_evaluations_surgeries;
DROP TABLE IF EXISTS dops_evaluations;
DROP TABLE IF EXISTS dops_forms;
DROP TABLE IF EXISTS evaluation_forms;
DROP TABLE IF EXISTS orthopedic_surgeries_activity_review_surgeries;
DROP TABLE IF EXISTS orthopedic_surgeries_activity_review;
DROP TABLE IF EXISTS orthopedic_surgery_activities_surgeries;
DROP TABLE IF EXISTS orthopedic_surgery_activities;
DROP TABLE IF EXISTS surgeries;
DROP TABLE IF EXISTS organizational_unit_specialties;
DROP TABLE IF EXISTS specialties_activity_types;
DROP TABLE IF EXISTS specialties;
DROP TABLE IF EXISTS surgery_diagnosis;
DROP TABLE IF EXISTS surgery_methods;
DROP TABLE IF EXISTS examinations_activities_reviews;
DROP TABLE IF EXISTS examinations_activity_reviews;
DROP TABLE IF EXISTS examination_activities;
DROP TABLE IF EXISTS examinations;
DROP TABLE IF EXISTS procedures;
DROP TABLE IF EXISTS practical_activity_types;
DROP TABLE IF EXISTS user_organizational_unit_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS organizational_units;
DROP TABLE IF EXISTS organizational_unit_types;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd






