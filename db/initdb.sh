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
    ('Alis', 'Williams', CURRENT_DATE(), 'Kanzas', 'travels, books', 'email1@examle.com', '$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse', CURRENT_TIME(), CURRENT_TIME()),
    ('Sophia', 'Johnson', CURRENT_DATE(), 'Chicago', 'skiing, games', 'email2@examle.com', '$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse', CURRENT_TIME(), CURRENT_TIME()),
    ('Emma', 'Brown', CURRENT_DATE(), 'Seattle', 'karaoke, music', 'email3@examle.com', '$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse', CURRENT_TIME(), CURRENT_TIME()),
    ('Isabella', 'Davis', CURRENT_DATE(), 'Kanzas', 'travels, reading', 'email4@examle.com', '$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse', CURRENT_TIME(), CURRENT_TIME()),
    ('Olivia', 'Miller', CURRENT_DATE(), 'Boston', 'karate, coloring', 'email5@examle.com', '$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse', CURRENT_TIME(), CURRENT_TIME()),
    ('Emily', 'Jones', CURRENT_DATE(), 'Los Angeles', 'makeup, books', 'email6@examle.com', '$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse', CURRENT_TIME(), CURRENT_TIME()),
    ('Katrin', 'Wilson', CURRENT_DATE(), 'San Francisco', 'magic, cooking', 'email7@examle.com', '$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse', CURRENT_TIME(), CURRENT_TIME());
EOSQL