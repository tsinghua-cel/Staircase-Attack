use mysql;

CREATE USER 'eth' @'127.0.0.1' identified BY '12345678';

CREATE USER 'eth' @'%' identified BY '12345678';

GRANT ALL ON *.* TO 'eth' @'%';

GRANT ALL ON *.* TO 'eth' @'127.0.0.1';

CREATE schema eth;
