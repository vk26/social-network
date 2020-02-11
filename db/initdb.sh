query="
  CREATE DATABASE IF NOT EXISTS social_dev;
  USE social_dev;
  
  DROP TABLE IF EXISTS users;

  CREATE TABLE users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name TEXT NOT NULL,
			surname TEXT NOT NULL,
			birthday DATE,
			city TEXT,
			about TEXT,
      avatar TEXT, 
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME,
      updated_at DATETIME
	) ENGINE INNODB;

	CREATE FULLTEXT INDEX fulltext_name_idx ON users(name, surname);
"
docker exec -it mysql_master sh -c "mysql -u root -p -e '$query'"