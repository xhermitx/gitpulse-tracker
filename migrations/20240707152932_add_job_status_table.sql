-- +goose Up
-- +goose StatementBegin
CREATE TABLE `job_status` (
  `job_id`        INT NOT NULL AUTO_INCREMENT,
  `triggered_at`    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `completed_at`  TIMESTAMP,
  PRIMARY KEY (`job_id`),
  FOREIGN KEY (`job_id`) REFERENCES `jobs` (`job_id`) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `jobs`;
-- +goose StatementEnd
