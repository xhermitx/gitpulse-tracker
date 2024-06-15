-- +goose Up
-- +goose StatementBegin
CREATE TABLE `candidates_list` (
  `job_id`    INT NOT NULL,
  `github_id` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`job_id`, `github_id`),
  FOREIGN KEY (`job_id`) REFERENCES `jobs` (`job_id`) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `candidates_list`;
-- +goose StatementEnd
