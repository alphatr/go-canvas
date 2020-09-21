package canvas

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// ImageCache 图片缓存
type ImageCache struct {
	Image  []byte    `json:"image"`
	Expire time.Time `json:"expire"`
}

// ImageCacheBucket 图片缓存
const ImageCacheBucket = "image-cache"

// BoltDB var
var BoltDB *bolt.DB

// Init 初始化 Bolt
func Init(filePath string) error {
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return err
	}

	BoltDB = db
	err = BoltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ImageCacheBucket))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// SaveCache 保存图片缓存
func SaveCache(key string, image []byte, expire time.Duration) error {
	encoded, err := json.Marshal(&ImageCache{image, time.Now().Add(expire)})
	if err != nil {
		return err
	}

	err = BoltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ImageCacheBucket))
		if err := bucket.Put([]byte(key), encoded); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// GetCache 获取图片缓存
func GetCache(key string) []byte {
	var dataByte []byte

	err := BoltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ImageCacheBucket))
		dataByte = bucket.Get([]byte(key))
		return nil
	})

	if err != nil {
		return nil
	}

	if len(dataByte) == 0 {
		return nil
	}

	cache := &ImageCache{}
	if err := json.Unmarshal(dataByte, cache); err != nil {
		return nil
	}

	if cache.Expire.Before(time.Now()) {
		return nil
	}

	clearExpiredCache()
	return cache.Image
}

func clearExpiredCache() error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ImageCacheBucket))

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			data := &ImageCache{}
			err := json.Unmarshal(value, data)

			if err != nil || data.Expire.Before(time.Now()) {
				if err := bucket.Delete(key); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
