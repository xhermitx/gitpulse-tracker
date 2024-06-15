-- +goose Up
-- +goose StatementBegin
CREATE TABLE `jobs` (
  `job_id`        INT NOT NULL AUTO_INCREMENT,
  `job_name`      VARCHAR(255) NOT NULL,
  `description`   VARCHAR(255),
  `drive_link`    VARCHAR(255),
  `created_at`    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `recruiter_id`  INT,
  PRIMARY KEY (`job_id`),
  FOREIGN KEY (`recruiter_id`) REFERENCES `recruiters` (`recruiter_id`) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `jobs`;
-- +goose StatementEnd
