INSERT INTO Roles (ID, Name, uid, rights) VALUES ('1', 'User', '1daefb8f-ed39-4942-b637-ae2b97fdbf01', '{"read":1}');
INSERT INTO Roles (ID, Name, uid, rights) VALUES ('2', 'Editor', '7cee1cdc-158d-4c84-a8d1-2ee6b7dd7d9e', '{"read":1,"write":1}');
INSERT INTO Roles (ID, Name, uid, rights) VALUES ('3', 'Admin', '73d59c8e-5df0-4a6d-a4d2-8b15f2f3cba8', '{"read":1,"write":1,"create":1}');
INSERT INTO Roles (ID, Name, uid, rights) VALUES ('4', 'Owner', '330ab705-f74b-4208-9a3c-1b91c64d5bf3', '{"read":1,"write":1,"create":1,"delete":1}');
SELECT setval(pg_get_serial_sequence('Roles', 'id'), coalesce(max(id), 0)+1 , false) FROM Roles;



INSERT INTO Users (ID, Name, uid, email, status, Role_Id) VALUES ('1', 'Test User', 'f11e8023-cb2c-41c2-b463-11b4880210b7', 'test@test.com', '1', '1');
INSERT INTO Users (ID, Name, uid, email, status, Role_Id) VALUES ('2', 'Developer 1', '5aacf1c4-e792-44b9-9bbf-accc3df4bb3d', 'dev1@test.com', '1', '2');
INSERT INTO Users (ID, Name, uid, email, status, Role_Id) VALUES ('3', 'Developer 2', '71f6ad63-8e92-41b8-9be0-084cfd8cf0cf', 'dev2@test.com', '1', '4');
SELECT setval(pg_get_serial_sequence('Users', 'id'), coalesce(max(id), 0)+1 , false) FROM Users;



INSERT INTO Groups (ID, Name, uid) VALUES ('1', 'Developers', '06cae884-c68c-418b-822a-7b9ec3fe3752');
INSERT INTO Groups (ID, Name, uid) VALUES ('2', 'Product Management', 'a3b01de0-5dd1-47b1-9c98-9616b1ee0639');
SELECT setval(pg_get_serial_sequence('Groups', 'id'), coalesce(max(id), 0)+1 , false) FROM Groups;



INSERT INTO Groups_Users_Map (Users_Id, Groups_Id) VALUES (1, 2);
INSERT INTO Groups_Users_Map (Users_Id, Groups_Id) VALUES (2, 1);
INSERT INTO Groups_Users_Map (Users_Id, Groups_Id) VALUES (3, 1);
INSERT INTO Groups_Users_Map (Users_Id, Groups_Id) VALUES (3, 2);