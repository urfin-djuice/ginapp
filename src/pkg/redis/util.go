package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func Ping() error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("cannot 'PING' db: %v", err)
	}
	return nil
}

func Get(key string) ([]byte, error) {
	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			return data, fmt.Errorf("error getting key %s: %v", key, err)
		}
	}
	return data, err
}

func Set(key string, value []byte) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}
func appendArgs(dst, src []interface{}) []interface{} {
	if len(src) == 1 {
		if ss, ok := src[0].([]string); ok {
			for _, s := range ss {
				dst = append(dst, s)
			}
			return dst
		}
	}
	dst = append(dst, src...)

	return dst
}

func SAdd(key string, members []interface{}) error {
	args := make([]interface{}, 2, 2+len(members))
	args[0] = key
	args = appendArgs(args, members)

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", args...)
	if err != nil {
		return fmt.Errorf("error setting set by key %s: %v", key, err)
	}
	return err
}

func SRandMember(key string) ([]byte, error) {
	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("SRANDMEMBER", key))
	if err != nil {
		return data, fmt.Errorf("error getting SRANDMEMBER %s: %v", key, err)
	}
	return data, err
}

func SetEx(key string, value []byte, seconds int32) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value, "EX", seconds)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting ex key %s to %s: %v", key, v, err)
	}
	return err
}

func Exists(key string) (bool, error) {
	conn := Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func Delete(key string) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func GetKeys(pattern string) ([]string, error) {
	conn := Pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

func Incr(counterKey string) (int, error) {
	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counterKey))
}
