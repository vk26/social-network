mysql -u "${DB_MYSQL_USER}" --password="${DB_MYSQL_PASSWORD}" <<-EOSQL
  CREATE DATABASE IF NOT EXISTS social_dev;
  USE social_dev;
  
  DROP TABLE IF EXISTS users;

  CREATE TABLE users (
			id INT AUTO_INCREMENT,
			name TEXT NOT NULL,
			surname TEXT NOT NULL,
			birthday DATE,
			city TEXT,
			about TEXT,
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME,
      updated_at DATETIME,
			PRIMARY KEY (id)
	);

  INSERT INTO users(name, surname, birthday, city, about, email, password_hash, created_at, updated_at) 
  VALUES
    ('Alis', 'Nash', CURRENT_DATE(), 'Kanzas', 'travels', 'alis@examle.com', 'asdf', CURRENT_TIME(), CURRENT_TIME()),
    ('Bob', 'Hans', CURRENT_DATE(), 'SF', 'games', 'bob@examle.com', 'aqwe', CURRENT_TIME(), CURRENT_TIME());
EOSQL