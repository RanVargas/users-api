CREATE EXTENSION "uuid-ossp";


DROP TABLE IF EXISTS Roles;
CREATE TABLE Roles (id SERIAL PRIMARY KEY, name VARCHAR NOT NULL, uid UUID, rights JSONB);
CREATE INDEX roles_uid_idx ON Roles USING BTREE (uid);



DROP TABLE IF EXISTS Users;
CREATE TABLE Users (id SERIAL PRIMARY KEY, name VARCHAR NOT NULL, uid UUID, email VARCHAR, status Integer, Role_id Integer);
CREATE INDEX users_uid_idx ON Users USING BTREE (uid);
CREATE INDEX users_role_id_idx ON Users USING BTREE (Role_id);



DROP TABLE IF EXISTS Groups;
CREATE TABLE Groups (id SERIAL PRIMARY KEY, name VARCHAR NOT NULL, uid UUID);
CREATE INDEX groups_uid_idx ON Groups USING BTREE (uid);



DROP TABLE IF EXISTS Groups_Users_Map;
CREATE TABLE Groups_Users_Map (id SERIAL PRIMARY KEY, Users_Id Integer, Groups_Id Integer);

