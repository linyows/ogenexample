CREATE TABLE pets (
  id BIGINT AUTO_INCREMENT NOT NULL,
  name VARCHAR(255) NOT NULL,
  `status` ENUM('available', 'pending', 'sold') NOT NULL,

  PRIMARY KEY (id)
);
