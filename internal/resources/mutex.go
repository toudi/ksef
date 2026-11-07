package resources

import "sync"

// this mutex is a consequence of the two independent things:
// - rewrite of invoicedb to interface-based approcah which changes the way how client is initialized
// - paralell download functionality
// there are some resources that are shared (i.e. certsdb, keyring and so on) and some that are
// initialized on a per-NIP basis. because they are initialized per-NIP, they need to live in a map.
// and because of that, we've got to protect it with a mutex since the map reallocates memory.
// so within the `client.go` and `invoicedb.go` you will see a pattern of public accessors that wrap
// mutex locking and unlocking and call corresponding private functions which end up with `*Locked`
var mu sync.Mutex
