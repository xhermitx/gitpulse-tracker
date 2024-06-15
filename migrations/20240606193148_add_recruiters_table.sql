-- +goose Up
-- +goose StatementBegin
CREATE TABLE `recruiters` (
  `recruiter_id`    INT NOT NULL AUTO_INCREMENT,
  `username`        VARCHAR(255) UNIQUE NOT NULL,
  `password`        VARCHAR(255) NOT NULL,
  `email`           VARCHAR(255) NOT NULL,
  `company`         VARCHAR(255) NOT NULL,
  `created_at`      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`recruiter_id`)
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `recruiters`;
-- +goose StatementEnd
