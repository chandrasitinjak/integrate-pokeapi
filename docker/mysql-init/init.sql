CREATE DATABASE IF NOT EXISTS pokemon;

USE pokemon;

CREATE TABLE IF NOT EXISTS pokemon (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255),
    rate INT DEFAULT 0,
    gender VARCHAR(50),
    PRIMARY KEY (id)
);

-- Tambahkan unique constraint untuk kombinasi name + gender
ALTER TABLE pokemon 
ADD CONSTRAINT uq_name_gender UNIQUE (name, gender);
