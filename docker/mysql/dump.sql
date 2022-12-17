/* Table to store capacity of tables*/
CREATE TABLE `tables` (
  `id` INT UNSIGNED NOT NULL auto_increment,
  `capacity` INT NOT NULL,
  PRIMARY KEY (`id`)
);

/* Table to stroe guest lists*/
CREATE TABLE `guests` (
  `id` INT UNSIGNED NOT NULL,
  `name` VARCHAR(100) NOT NULL UNIQUE,
  `accompanying_guests` INT UNSIGNED,
  `status` ENUM('allotted', 'checked-in', 'checked-out') DEFAULT 'allotted',
  `time_arrived` TIMESTAMP,
  FOREIGN KEY (`id`) REFERENCES tables(`id`)
);
