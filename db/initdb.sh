mysql -u "${DB_MYSQL_USER}" --password="${DB_MYSQL_PASSWORD}" <<-EOSQL
  CREATE DATABASE IF NOT EXISTS social_dev;
  USE social_dev;
  
  DROP TABLE IF EXISTS users;

  CREATE TABLE users (
			id INT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			surname TEXT NOT NULL,
			birthday DATE,
			city TEXT,
			about TEXT,
      avatar TEXT, 
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME,
      updated_at DATETIME,
	) ENGINE INNODB;

  CREATE INDEX name_idx ON users(name(15));
  CREATE INDEX surname_idx ON users(surname(15));
EOSQL