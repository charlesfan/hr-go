package db

// admin user, pwd = admin@init
var InitSqlStr = `INSERT IGNORE INTO hr_employee VALUES ('60aacf8b-4ff9-4d46-bd53-f14c0fbf9a8e','admin','admin@init.com','$2a$10$C1PYnbqMWn3POHSf/cOlBOd0/kSfE89KRhisk1e5yujl9rwTPTAoG','hr',UNIX_TIMESTAMP(NOW()),NOW(),NOW(),NULL);`
