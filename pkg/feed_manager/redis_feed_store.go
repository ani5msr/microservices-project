package feed_manager

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/go-redis/redis"
	"github.com/pelletier/go-toml"
)

const redisMaxPageSize = 10

// RedisFeedStore manages a UserEvents data structure
type RedisFeedStore struct {
	redis *redis.Client
}

func (m *RedisFeedStore) GetFeed(username string, startIndex int) (events []*om.PostManagerEvent, nextIndex int, err error) {
	stop := startIndex + redisMaxPageSize - 1
	result, err := m.redis.LRange(username, int64(startIndex), int64(stop)).Result()
	if err != nil {
		return
	}

	for _, t := range result {
		var event om.PostManagerEvent
		err = toml.Unmarshal([]byte(t), &event)
		if err != nil {
			return
		}

		events = append(events, &event)
	}

	if len(result) == redisMaxPageSize {
		nextIndex = stop + 1
	} else {
		nextIndex = -1
	}

	return
}

func (m *RedisFeedStore) AddEvent(username string, event *om.PostManagerEvent) (err error) {
	t, err := toml.Marshal(*event)
	if err != nil {
		return
	}

	err = m.redis.RPush(username, t).Err()
	return
}

func NewRedisFeedStore(address string) (store Store, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // use empty password for simplicity. should come from a secret in production
		DB:       0,  // use default DB
	})

	_, err = client.Ping().Result()
	if err != nil {
		return
	}

	store = &RedisFeedStore{redis: client}
	return
}
