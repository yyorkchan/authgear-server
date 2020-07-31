package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	goredis "github.com/gomodule/redigo/redis"

	"github.com/authgear/authgear-server/pkg/auth/config"
	"github.com/authgear/authgear-server/pkg/redis"
)

type StoreRedis struct {
	Redis  *redis.Handle
	AppID  config.AppID
	Config *config.VerificationConfig
}

func (s *StoreRedis) Create(code *Code) error {
	data, err := json.Marshal(code)
	if err != nil {
		return err
	}

	return s.Redis.WithConn(func(conn redis.Conn) error {
		codeKey := redisCodeKey(s.AppID, code.Code)
		ttl := toMilliseconds(s.Config.CodeExpiry.Duration())
		_, err := goredis.String(conn.Do("SET", codeKey, data, "PX", ttl, "NX"))
		if errors.Is(err, goredis.ErrNil) {
			return errors.New("duplicated code")
		} else if err != nil {
			return err
		}

		return nil
	})
}

func (s *StoreRedis) Get(code string) (*Code, error) {
	key := redisCodeKey(s.AppID, code)
	var codeModel *Code
	err := s.Redis.WithConn(func(conn redis.Conn) error {
		data, err := goredis.Bytes(conn.Do("GET", key))
		if errors.Is(err, goredis.ErrNil) {
			return ErrCodeNotFound
		} else if err != nil {
			return err
		}

		err = json.Unmarshal(data, &codeModel)
		if err != nil {
			return err
		}

		return nil
	})
	return codeModel, err
}

func (s *StoreRedis) Delete(code string) error {
	return s.Redis.WithConn(func(conn redis.Conn) error {
		key := redisCodeKey(s.AppID, code)
		_, err := conn.Do("DEL", key)
		if err != nil {
			return err
		}
		return err
	})
}

func toMilliseconds(d time.Duration) int64 {
	return int64(d / time.Millisecond)
}

func redisCodeKey(appID config.AppID, code string) string {
	return fmt.Sprintf("app:%s:verification-code:%s", appID, code)
}
