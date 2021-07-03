package memkv

import (
	"fmt"
	"github.com/tidwall/gjson"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type Timer struct {
	begin time.Time
	end   time.Time
}

func NewTimer() *Timer {
	t := &Timer{}
	t.Start()
	return t
}
func (t *Timer) Start() {
	t.begin = time.Now()
}

func (t *Timer) Stop() float64 {
	elapsed := time.Since(t.begin)
	return elapsed.Seconds() * 1000
}

func TestIndexString(t *testing.T) {

	//var item *DBItem
	//item = &DBItem{}

}

type Dimension struct {
	Dim1  uint64
	Dim2  string
	Dim3  string
	Dim4  string
	Dim5  string
	Value float64
}

func IndexJSON2(path string) func(a, b string) bool {
	return func(a, b string) bool {
		v1 := gjson.Get(a, path)
		v2 := gjson.Get(b, path)
		if v1.Str == v2.Str {
			return true
		} else {
			return false
		}
		r := gjson.Get(a, path).Less(gjson.Get(b, path), true)
		return r
	}
}

var gDB *DB

func TestDB_Indexes01(t *testing.T) {
	db, _ := Open(":memory:")
	db.CreateIndex("RawKeyAndCommitTSIndex", "*", IndexRawKey, IndexCommitTS)
	db.CreateIndex("KeyIndex", "*", IndexKey)
	var item *DBItem
	item = &DBItem{Key: createKey("B", 5, 8), RawKey: []byte("A1"), StartTS: 5, CommitTS: 80}
	db.Put(item)
	item = &DBItem{Key: createKey("B", 2, 5), RawKey: []byte("A2"), StartTS: 5, CommitTS: 70}
	db.Put(item)
	item = &DBItem{Key: createKey("B", 3, 6), RawKey: []byte("A3"), StartTS: 6, CommitTS: 40}
	db.Put(item)
	item = &DBItem{Key: createKey("A", 4, 7), RawKey: []byte("A4"), StartTS: 7, CommitTS: 70}
	db.Put(item)
	item = &DBItem{Key: createKey("A", 5, 4), RawKey: []byte("A5"), StartTS: 4, CommitTS: 40}
	db.Put(item)

	db.View(func(tx *Tx) error {
		tx.Ascend("KeyIndex", func(key Key, value *DBItem) bool {
			fmt.Println(*value)
			return true
		})
		return nil
	})
}
func createKey(prefix string, index int, data uint64) []byte {
	key := prefix + strconv.Itoa(index)
	b := []byte(key)
	b = mvccEncode(b, data)
	return b
}
func TestDB_Indexes04(t *testing.T) {
	//testLoadData(t)
	//var err error
	//timer := NewTimer()
	////var result string
	//const NUM = 1
	////for i :=1;i <= NUM;i++{
	//val := `{"Dim2":"AAA1","Dim3":"BBB1"}`
	//err = gDB.View(func(tx *Tx) error {
	//	err := tx.AscendEqual("Dim2_Dim3", val, func(key, value string) bool {
	//
	//		fmt.Printf("%s: %s\n", key, value)
	//
	//		return true
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	return nil
	//})
	//if err != nil {
	//	fmt.Println(err)
	//}
	////}
	//fmt.Println(timer.Stop())

}
func TestTree_Items(t *testing.T) {
	//list := items{}
	//a := &DBItem{}
	//b := &DBItem{}
	//list.addItems(a)
	//list.addItems(b)
	//fmt.Print(list.len())

}
func TestDB_Indexes05(t *testing.T) {
	db, _ := Open(":memory:")
	db.CreateIndex("IndexRawKey", "*", IndexRawKey)
	db.CreateIndex("KeyIndex", "*", IndexKey)
	for i := 1; i <= 2000; i++ {
		str := strconv.Itoa(i)
		db.Put(&DBItem{Key: []byte(str), RawKey: []byte(str), Val: []byte(str)})
		db.Delete([]byte(str))
	}
	db.Ascend("IndexRawKey", func(key Key, value *DBItem) bool {
		fmt.Println(string(key))
		return true
	})
}

func TestDBItem_MarshalBinary(t *testing.T) {
	dbItem := &DBItem{}
	dbItem.Key = []byte{1, 2, 3}
	dbItem.RawKey = []byte{1, 2, 3, 4}
	dbItem.Val = []byte{1, 2, 3, 4, 5}
	dbItem.ValueType = 1
	dbItem.StartTS = 2
	dbItem.CommitTS = 3
	dbItem.Op = 4
	dbItem.Ttl = 5
	dbItem.ForUpdateTS = 6
	dbItem.TxnSize = 7
	dbItem.MinCommitTS = 8
	buff, _ := dbItem.MarshalBinary()
	dbItem2 := &DBItem{}
	dbItem2.UnmarshalBinary(buff)
	if !reflect.DeepEqual(dbItem, dbItem2) {
		t.Error("序列化失败!")
	}
}
func TestCreatDB2(t *testing.T) {
	length := 5
	fmt.Println(1 << length)
}
func TestCreatDB(t *testing.T) {
	db, err := CreatDB(1, true)
	if err != nil {
		t.Error(err)
	}
	dbItem := &DBItem{}
	dbItem.Key = []byte{1, 2, 3}
	dbItem.RawKey = []byte{1, 2, 3, 4}
	dbItem.Val = []byte{1, 2, 3, 4, 5}
	dbItem.ValueType = 1
	dbItem.StartTS = 2
	dbItem.CommitTS = 3
	dbItem.Op = 4
	dbItem.Ttl = 5
	dbItem.ForUpdateTS = 6
	dbItem.TxnSize = 7
	dbItem.MinCommitTS = 8
	db.Put(dbItem)
	db.Close()

	db2, err := CreatDB(1, true)
	if err != nil {
		t.Error(err)
	}
	//time.Sleep(3 * time.Second)
	db2.LoadDB()
	key := []byte{1, 2, 3}
	dbItem2 := db2.Get(key)

	if !reflect.DeepEqual(dbItem, dbItem2) {
		t.Error("数据获取失败!")
	}
}
