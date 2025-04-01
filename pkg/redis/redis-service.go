package redisdb

import "log"

// ========================= Redis String Operations =========================

func StringGet(db int, key string) string {
	rdb := GetRedisClient(db)
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Redis Get Error (DB %d): %v", db, err)
		return ""
	}
	return val
}

func StringSet(db int, key string, value string) {
	rdb := GetRedisClient(db)
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Redis Set Error (DB %d): %v", db, err)
	}
}

// ========================= Redis List Operations =========================

func ListSet(db int, key string, value string) {
	rdb := GetRedisClient(db)
	err := rdb.LPush(ctx, key, value).Err()
	if err != nil {
		log.Printf("Redis LPush Error (DB %d): %v", db, err)
	}
}

func ListDelElement(db int, key string, value string) {
	rdb := GetRedisClient(db)
	_, err := rdb.LRem(ctx, key, 1, value).Result()
	if err != nil {
		log.Printf("Error removing element from list (DB %d): %v", db, err)
	}
}
func ListMove(db int, source, destination, value string) {
	ListDelElement(db, source, value)
	ListSet(db, destination, value)
}

// ========================= Redis Hash Operations =========================

func HashSet(db int, key, field, value string) {
	rdb := GetRedisClient(db)
	err := rdb.HSet(ctx, key, field, value).Err()
	if err != nil {
		log.Printf("Redis HSet Error (DB %d): %v", db, err)
	}
}

func HashGet(db int, key, field string) string {
	rdb := GetRedisClient(db)
	val, err := rdb.HGet(ctx, key, field).Result()
	if err != nil {
		log.Printf("Redis HGet Error (DB %d): %v", db, err)
		return ""
	}
	return val
}

func HashDelKey(db int, key, field string) {
	rdb := GetRedisClient(db)
	err := rdb.HDel(ctx, key, field).Err()
	if err != nil {
		log.Printf("Redis HDel Error (DB %d): %v", db, err)
	}
}

func HashDelAll(db int, key string) {
	rdb := GetRedisClient(db)
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Redis Del Error (DB %d): %v", db, err)
	}
}

func HashGetAll(db int, key string) map[string]string {
	rdb := GetRedisClient(db)
	val, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// ========================= General Redis Operations =========================

func GetAllList(db int) (map[string][]string, error) {
	rdb := GetRedisClient(db)
	cursor := uint64(0)
	listData := make(map[string][]string) // To hold the list keys and their values

	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			log.Printf("Redis Scan Error: %v", err)
			return nil, err
		}

		for _, key := range keys {
			keyType, err := rdb.Type(ctx, key).Result()
			if err != nil {
				log.Printf("Redis Type Check Error: %v", err)
				return nil, err
			}

			if keyType == "list" {
				values, err := rdb.LRange(ctx, key, 0, -1).Result()
				if err != nil {
					log.Printf("Redis LRange Error: %v", err)
					return nil, err
				}
				listData[key] = values
			}
		}

		if cursor == 0 { // No more keys left
			break
		}
	}
	return listData, nil
}
