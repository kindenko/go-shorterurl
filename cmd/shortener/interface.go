package main

func ee() {

}

// package main
// //package storage, file storage.go
// type UrlStorer interface {
// 	Save(short, long string ) error
// 	Get(short string) (string, error)
// }

// //package storage, file memstorage.go
// type MemoryStorage struct {
// 	store map[string]string
// }

// func (m *MemoryStorage) Save(short, long string) error {
// 	m.store[short] = long
// 	return nil
// }

// func (m *MemoryStorage) Get(short string) (string, error) {
// 	long, ok := m.store[short]
// 	if ok {
// 		return long, nil
// 	}
// 	return "", error("steringfdsgfdvfdshgskjdh")
// }

// //dbstorage.go
// //filestorage.go

// func InitStorage(cfg Config) UrlStorer {
// 	if (cfg.DBSN == "") {
// 		return MemoryStorage{
// 			store == map[string]string{}
// 		}
// 	}
// }

// type B struct{}
// func (b *B) Save()

// func DoSomething(a UrlStorer) {
// 	a.Save()
// }
