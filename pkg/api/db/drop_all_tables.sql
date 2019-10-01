drop table if exists team cascade;
drop table if exists users cascade;
drop table if exists application cascade;
drop table if exists package cascade;
drop table if exists flatcar_action cascade;
drop table if exists channel cascade;
drop table if exists groups cascade;
drop table if exists instance cascade;
drop table if exists instance_status cascade;
drop table if exists instance_application cascade;
drop table if exists instance_status_history cascade;
drop table if exists event_type cascade;
drop table if exists event cascade;
drop table if exists activity cascade;
drop table if exists package_channel_blacklist cascade;
drop table if exists database_migrations;
-- Legacy tables if we're dropping tables in a non-migrated DB
drop table if exists coreos_action cascade;