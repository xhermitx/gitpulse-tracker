package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	redis "github.com/redis/go-redis/v9"
	"github.com/xhermitx/gitpulse-tracker/profiler-service/models"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(rdb *redis.Client) *RedisClient {
	return &RedisClient{
		client: rdb,
	}
}

func (rdb RedisClient) Set(ctx context.Context, profile models.TopCandidates) error {

	totalScore := float64(profile.Followers + profile.RepoStars)

	profile.Score = uint(totalScore)

	// Use the JobID to create a unique key for each job's Sorted Set
	key := fmt.Sprintf("job:%d:top_candidates", profile.JobId)

	profileJson, err := json.Marshal(profile)
	if err != nil {
		log.Println(err)
	}

	log.Println("Profile JSON Length: ", len(profileJson))

	// Add the candidate to the Sorted Set with the totalScore as the score
	_, err = rdb.client.ZAdd(ctx, key, redis.Z{
		Score:  totalScore,
		Member: profileJson,
	}).Result()

	if err != nil {
		return err
	}

	fmt.Println("Successfully stored the candidate : ", profile.GithubId)

	// Only keep the top 5 candidates with the highest scores
	_, err = rdb.client.ZRemRangeByRank(ctx, key, 0, -6).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rdb RedisClient) Get(ctx context.Context, jobID uint) ([]models.TopCandidates, error) {

	key := fmt.Sprintf("job:%d:top_candidates", jobID)

	res, err := rdb.client.ZRevRangeWithScores(ctx, key, 0, 4).Result()
	if err != nil {
		return nil, err
	}

	var topCandidates []models.TopCandidates

	for _, z := range res {
		var candidate models.TopCandidates

		// Since z.Member is an interface{}, assert it as byte slice
		candidateJson, ok := z.Member.(string)
		if !ok {
			return nil, fmt.Errorf("member is not a byte slice")
		}

		err := json.Unmarshal([]byte(candidateJson), &candidate)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal member: %w", err)
		}

		fmt.Printf("Member: %s, Score: %f\n", candidate.GithubId, z.Score)
		topCandidates = append(topCandidates, candidate)
	}

	return topCandidates, nil
}
