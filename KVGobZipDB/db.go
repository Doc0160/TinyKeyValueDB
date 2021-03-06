package KVGobZipDB

import(
    "compress/gzip"
	"encoding/gob"
	"io/ioutil"
	"bytes"
	"sync"
	"errors"
    "fmt"
)


var (
	ErrNotFound = errors.New("db: key not found")
	ErrBadValue = errors.New("db: bad value")
    //ErrAlreadyExist = errors.New("db: key already exist")
)

type DBData struct{
	Keys []string
	Values []string
}

type DB struct{
	data DBData
	sync.Mutex
	filename string
}

func Open(filename string) DB {
	db := DB{
        DBData{},
        sync.Mutex{},
        filename,
    }
    //db.Lock()
    //defer db.Unlock()
	b, err := ioutil.ReadFile(filename)
	if err == nil {
        br := bytes.NewReader(b)
        r, err := gzip.NewReader(br)
        if err == nil {
            err := gob.NewDecoder(r).Decode(&db.data)
            r.Close()
            if err != nil {
                fmt.Println(err)
            }
        }else{
            fmt.Println(err)
        }
	} else {
		fmt.Println(err)
	}
	return db
}

/*func (db*DB)Key(i int, s *string){
    //db.Lock()
    *s =  db.data.Keys[i]
    //db.Unlock()
}

func (db*DB)Value(i int, r interface{}){
    //db.Lock()
    r = db.data.Values[i]
    //db.Unlock()
}*/

func (db*DB) Save() {
	var buf bytes.Buffer
    
    db.Lock()
	if err := gob.NewEncoder(&buf).Encode(db.data); err != nil {
		fmt.Println(err)
	}
    db.Unlock()
    
	var b bytes.Buffer
    w, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
    w.Write(buf.Bytes())
    w.Close()
	ioutil.WriteFile(db.filename, b.Bytes(), 0777)
}

func (db*DB) Put(key string, value interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return nil
	}
    
    db.Lock()
	i := index(db.data.Keys, key)
	if i == -1 {
		db.data.Keys = append(db.data.Keys, key)
		db.data.Values = append(db.data.Values, string(buf.Bytes()))
	}else{
		db.data.Values[i] = string(buf.Bytes())
	}
    db.Unlock()
    
	return nil
}

func (db*DB) Get(key string, value interface{}) error {
    var err = ErrNotFound
    
    db.Lock()
	i := index(db.data.Keys, key)
	if i != -1 {
		err = gob.NewDecoder(bytes.NewReader([]byte(db.data.Values[i]))).Decode(value)
	}
    db.Unlock()
    
	return err
}

func (db*DB) Delete(key string) error {
	db.Lock()
	i := index(db.data.Keys, key)
	if i == -1 {
        db.Unlock()
		return ErrNotFound
	}
	db.data.Keys = delete(db.data.Keys, i)
	db.data.Values = delete(db.data.Values, i)
    db.Unlock()
	return nil
}

func delete(a []string, i int) []string {
    copy(a[i:], a[i+1:])
    a[len(a)-1] = ""
    a = a[:len(a)-1]
	a = append(a[:i], a[i+1:]...)
	return a
}

func index(s []string, e string) int {
	for k, a := range s {
		if a == e {
			return k
		}
	}
    return -1
}
