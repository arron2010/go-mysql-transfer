package memkv

import (
	"github.com/xp/shorttext-db/bbolt"
	"github.com/xp/shorttext-db/gopool"
	"os"
	"time"
)

var fileMode os.FileMode = 0600
var (
	DB_NAME = []byte("memkv")
)

type ItemHandler func(item *DBItem)

type dbWriter struct {
	db   *bbolt.DB
	pool *gopool.PoolWithFunc
	sync bool
}

func newDbWriter(path string) (w *dbWriter, err error) {
	w = &dbWriter{}
	options := gopool.Options{}
	options.ExpiryDuration = time.Duration(3) * time.Second
	options.Nonblocking = false
	options.PreAlloc = true
	w.pool, _ = gopool.NewPoolWithFunc(10, func(item interface{}) {
		dbItem, ok := item.(*DBItem)
		if !ok {
			return
		}
		w.innerPut(dbItem)
		//logger.Infof("异步执行KV磁盘存储 key=%v\n",dbItem.Key)
	}, gopool.WithOptions(options))

	w.db, err = bbolt.Open(path, fileMode, &bbolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}
	err = w.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(DB_NAME)
		return err
	})
	return w, err
}

func (w *dbWriter) Close() (err error) {
	return w.db.Close()
}

func (w *dbWriter) Put(item *DBItem) {
	if w.sync {
		w.innerPut(item)
		//	logger.Infof("同步执行KV磁盘存储 key=%v\n",item.Key)

	} else {
		w.pool.Invoke(item)
	}
}
func (w *dbWriter) innerPut(item *DBItem) {
	buf, err := item.MarshalBinary()
	if err != nil {
		logger.Errorf("key=%v写入bbolt db失败\n", item.Key)
		return
	}
	w.put(item.Key, buf)

}

func (w *dbWriter) ForEach(handler ItemHandler) {
	w.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(DB_NAME)
		b.ForEach(func(k, v []byte) error {
			dstv := make([]byte, len(v))
			copy(dstv, v)
			item := &DBItem{}
			err := item.UnmarshalBinary(dstv)
			if err == nil {
				handler(item)
			} else {
				logger.Errorf("加载数据库发生错误,key=%v\n", k)
			}
			return nil
		})
		return nil
	})
}
func (w *dbWriter) get(key []byte) ([]byte, error) {
	var value []byte
	var err error
	err = w.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(DB_NAME)
		v := b.Get(key)
		if v != nil {
			value = append(value, b.Get(key)...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, err
}
func (w *dbWriter) put(key []byte, value []byte) error {
	err := w.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(DB_NAME)
		err := b.Put(key, value)
		return err
	})
	return err
}
