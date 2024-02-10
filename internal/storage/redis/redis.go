package redis

import (
	"bot/internal/config"
	"bot/pkg/logging"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
)

type Redis struct {
	log         *logging.Logger
	redisClient *redis.Client
}

func NewRedis(cfg *config.RedisDBConfig, log *logging.Logger) *Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       cfg.DB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("failed to connect redis", zap.Error(err))
	}

	return &Redis{
		log:         log,
		redisClient: redisClient,
	}
}

func (r *Redis) SetState(telegramID int64, stateName string, stateData *map[string]interface{}) error {
	stateJSON, err := json.Marshal(stateData)
	if err != nil {
		r.log.Error("failed to marshal state data", zap.Error(err))
		return err
	}
	err = r.redisClient.HSet(context.TODO(), strconv.FormatInt(telegramID, 10), stateName, stateJSON).Err()
	if err != nil {
		r.log.Error("failed to set state", zap.Error(err))
		return err
	}

	return nil
}

func (r *Redis) GetStateData(telegramID int64, stateName string) (map[string]interface{}, error) {
	stateJSON, err := r.redisClient.HGet(context.TODO(), strconv.FormatInt(telegramID, 10), stateName).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		r.log.Error("failed to get state", zap.Error(err))
		return nil, err
	}

	var stateData map[string]interface{}
	err = json.Unmarshal([]byte(stateJSON), &stateData)
	if err != nil {
		r.log.Error("failed to unmarshal state data", zap.Error(err))
		return nil, err
	}

	return stateData, nil
}

func (r *Redis) GetStateAll(telegramID int64) (map[string]interface{}, error) {
	stateValues, err := r.redisClient.HGetAll(context.TODO(), strconv.FormatInt(telegramID, 10)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		r.log.Error("failed to get state", zap.Error(err))
		return nil, err
	}

	stateData := make(map[string]interface{})
	for key, value := range stateValues {
		var data interface{}
		if err = json.Unmarshal([]byte(value), &data); err != nil {
			r.log.Error("failed to unmarshal state data", zap.Error(err))
			return nil, err
		}
		stateData[key] = data
	}

	return stateData, nil
}

func (r *Redis) UpdateStateMap(telegramID int64, stateName string, stateData map[string]interface{}) error {
	currentState, err := r.GetStateData(telegramID, stateName)
	if currentState == nil {
		_ = r.SetState(telegramID, stateName, &map[string]interface{}{"start": true})
		currentState, err = r.GetStateData(telegramID, stateName)
	}

	for key, value := range stateData {
		currentState[key] = value
	}

	err = r.SetState(telegramID, stateName, &currentState)
	if err != nil {
		r.log.Error("failed to update state", zap.Error(err))
		return err
	}
	return nil
}

func (r *Redis) UpdateState(telegramID int64, stateName string, fieldName string, fieldValue interface{}) error {
	currentState, err := r.GetStateData(telegramID, stateName)
	if currentState == nil {
		_ = r.SetState(telegramID, stateName, &map[string]interface{}{"start": true})
		currentState, err = r.GetStateData(telegramID, stateName)
	}

	currentState[fieldName] = fieldValue

	err = r.SetState(telegramID, stateName, &currentState)
	if err != nil {
		r.log.Error("failed to update state", zap.Error(err))
		return err
	}

	return nil
}

func (r *Redis) ClearState(telegramID int64, stateName string) error {
	err := r.redisClient.HDel(context.Background(), strconv.FormatInt(telegramID, 10), stateName).Err()
	if err != nil {
		r.log.Error("failed to clear state", zap.Error(err))
		return err
	}

	return nil
}

func (r *Redis) ClearAllStates(telegramID int64) error {
	err := r.redisClient.Del(context.Background(), strconv.FormatInt(telegramID, 10)).Err()
	if err != nil {
		r.log.Error("failed to clear all states", zap.Error(err))
		return err
	}

	return nil
}
