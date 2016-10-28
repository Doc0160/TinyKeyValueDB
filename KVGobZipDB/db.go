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
	//ErrBadValue = errors.New("db: bad value")
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
	db := DB{}
	db.filename = filename
	b, err := ioutil.ReadFile(filename)
	if err == nil {
        br := bytes.NewReader(b)
        r, err := gzip.NewReader(br)
        defer r.Close()
        if err == nil {
            db.Lock()
            err := gob.NewDecoder(r).Decode(&db.data)
            db.Unlock()
            if err != nil {
                fmt.Println(err)
            }
        }else{
            fmt.Println(err)
        }
	} else {
		fmt.Println(err)
		db.Save()
	}
	return db
}

func (db*DB) Select(fn func(string)bool) []string {
    db.Lock()
    defer db.Unlock()
    var s []string
    for _, k := range db.data.Keys {
        if fn(k) {
            s = append(s, k)
        }
    }
    return s
}

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
	db.Lock()
	defer db.Unlock()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return nil
	}
	i := index(db.data.Keys, key)
	if i == -1 {
		db.data.Keys = append(db.data.Keys, key)
		db.data.Values = append(db.data.Values, string(buf.Bytes()))
	}else{
		db.data.Values[i] = string(buf.Bytes())
	}
	return nil
}


func (db*DB) Get(key string, value interface{}) error {
	db.Lock()
	defer db.Unlock()
	i := index(db.data.Keys, key)
	if i != -1 {
		d := gob.NewDecoder(bytes.NewReader([]byte(db.data.Values[i])))
		return d.Decode(value)
	}
	return ErrNotFound
}

func (db*DB) Delete(key string) error {
	db.Lock()
	defer db.Unlock()
	i := index(db.data.Keys, key)
	if i == -1 {
		return ErrNotFound
	}
	db.data.Keys = delete(db.data.Keys, i)
	db.data.Values = delete(db.data.Values, i)
	return nil
}

func delete(a []string, i int) []string {
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
