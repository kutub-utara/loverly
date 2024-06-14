BEGIN;

INSERT INTO users (id, email, password, created_at, updated_at) VALUES 
(1, 'handsome@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(2, 'lady1@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(3, 'lady2@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(4, 'lady3@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(5, 'lady4@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(6, 'lady5@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(7, 'lady6@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(8, 'lady7@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(9, 'lady8@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(10, 'lady9@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(11, 'lady10@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(12, 'lady11@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now()),
(13, 'lady12@gmail.com', '$2a$12$o5tYZ1aW3TsTtBl1oYQdJ.VeyFw3trpqCfHI0aQJtL2UcQFfJtMOS', now(), now());

INSERT INTO profiles (user_id, name, gender, birthday, created_at, updated_at) VALUES
(1, 'Handsome', 'male', '1990-01-01', now(), now()),
(2, 'Lady 1', 'female', '1990-02-01', now(), now()),
(3, 'Lady 2', 'female', '1990-02-01', now(), now()),
(4, 'Lady 3', 'female', '1990-02-01', now(), now()),
(5, 'Lady 4', 'female', '1990-02-01', now(), now()),
(6, 'Lady 5', 'female', '1990-02-01', now(), now()),
(7, 'Lady 6', 'female', '1990-02-01', now(), now()),
(8, 'Lady 7', 'female', '1990-02-01', now(), now()),
(9, 'Lady 8', 'female', '1990-02-01', now(), now()),
(10, 'Lady 9', 'female', '1990-02-01', now(), now()),
(11, 'Lady 10', 'female', '1990-02-01', now(), now()),
(12, 'Lady 11', 'female', '1990-02-01', now(), now()),
(13, 'Lady 12', 'female', '1990-02-01', now(), now());

COMMIT;