-- +goose Up
-- +goose StatementBegin
CREATE TABLE `top_candidates` (
  `candidate_id` INT NOT NULL AUTO_INCREMENT,
  `github_id` VARCHAR(255) NOT NULL,
  `followers` INT,
  `contributions` INT,
  `most_popular_repo` VARCHAR(255),
  `repo_stars` INT,
  `score` INT,
  `job_id` INT,
  PRIMARY KEY (`candidate_id`),
  FOREIGN KEY (`job_id`) REFERENCES `jobs` (`job_id`) ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `top_candidates`;
-- +goose StatementEnd
