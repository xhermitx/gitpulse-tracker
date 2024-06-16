package api

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/models"
)

func Set(profile models.RedisCandidate, rdb *redis.Client, ctx context.Context) error {

	totalScore := float64(profile.Followers + profile.RepoStars)
	// Use the JobID to create a unique key for each job's Sorted Set
	key := fmt.Sprintf("job:%d:top_candidates", profile.JobId)

	// Add the candidate to the Sorted Set with the totalScore as the score
	_, err := rdb.ZAdd(ctx, key, redis.Z{
		Score:  totalScore,
		Member: profile,
	}).Result()

	if err != nil {
		return err
	}

	fmt.Println("Successfully stored the candidate : ", profile.GithubId)

	// Only keep the top 5 candidates
	rdb.ZRemRangeByRank(ctx, key, 0, -6)

	return nil
}

func Get(jobID uint, rdb *redis.Client, ctx context.Context) ([]models.RedisCandidate, error) {

	key := fmt.Sprintf("job:%d:top_candidates", jobID)

	res, err := rdb.ZRevRangeWithScores(ctx, key, 0, 4).Result()
	if err != nil {
		return nil, err
	}

	var topCandidates []models.RedisCandidate

	for _, z := range res {
		fmt.Printf("Member: %s, Score: %f\n", z.Member, z.Score)
		topCandidates = append(topCandidates, z.Member.(models.RedisCandidate))
	}

	return topCandidates, nil
}
