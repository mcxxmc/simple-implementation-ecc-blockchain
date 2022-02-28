package datacenter

import (
	"encoding/json"
	"errors"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/blockchain"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/client"
	"github.com/mcxxmc/simple-implementation-ecc/ecc"
	"github.com/mcxxmc/simple-implementation-ecc/galois"
)

type Record [512]byte	// assumes a record to be a byte array of 512 bytes.

// DataCenter the datacenter struct plays as the data center, which contains a blockchain holding past activity and
// records the public key of all the active users	//todo: incorporate a database to store both offline and online users
//
// Please use NewDataCenter() as the constructor, and call RandomInitialization() to set a new private key and the public key.
// Then you can use GetPublicKey() to get the corresponding public key.
type DataCenter struct {
	Self 					*client.Client				// the own ECDH instance of the datacenter wrapped by a Client object
	Bc 						*blockchain.Blockchain		// the blockchain
	Users					map[int]*galois.Point		// id: public key
	ActiveRecord	 		*Buffer						// the current record; will be stored in a new block when full
}

// NewDataCenter returns the pointer to a new DataCenter object.
//
// The elliptic curve must have been initialized with a generator point. (by ep.SetGeneratorPoint())
func NewDataCenter(ep *ecc.Elliptic, name string, id int) (*DataCenter, error) {
	bc, err := blockchain.NewBlockchain()
	if err != nil {
		return nil, err
	}
	dc := &DataCenter{
		Self: 			 client.NewClient(ep, name, id),
		Bc:              bc,
		Users:           make(map[int]*galois.Point),
		ActiveRecord: 	 NewBuffer(BufferSize),
	}
	return dc, nil
}

// RandomInitialization randomly sets a new private key and the corresponding public key.
//
// This can be used to update the keys as well.
func (dc *DataCenter) RandomInitialization() {
	dc.Self.RandomInitialization()
	dc.Users[dc.Self.Id] = dc.Self.PubKey		// add own pub key to the users
}

// RequestPublicKey returns a copy of the public key given the user id.
func (dc *DataCenter) RequestPublicKey(userId int)	(galois.Point, error) {
	if p, exist := dc.Users[userId]; !exist {
		return galois.NonePoint(), errors.New("user id not found")
	} else {
		return galois.NewPoint(p.X, p.Y), nil
	}
}

// RequestRegisterNewUser registers a new user.
func (dc *DataCenter) RequestRegisterNewUser(userId int, userPubKey galois.Point) error {
	if _, exist := dc.Users[userId]; exist {
		return errors.New("user id already exists")
	}
	dc.Users[userId] = &userPubKey
	return nil
}

// RequestUpdateUserPubKey updates the public key of a certain user; should verify the identity of the user first.
func (dc *DataCenter) RequestUpdateUserPubKey() {


}

// Write writes bytes of data into the current record, and creates new block as needed.
func (dc *DataCenter) Write(data interface{}) error {		// todo: add encryption (e.g., using gcm)
	if data == nil {
		return errors.New("no data to write")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	size := len(bytes)
	if size == 0 {
		return errors.New("data size 0")
	}

	pointer := 0	// points to the bytes
	for pointer < size {
		pointer += dc.ActiveRecord.ReadTillFull(bytes[pointer:])
		if dc.ActiveRecord.IsFull() {
			content := dc.ActiveRecord.Bytes()
			dc.ActiveRecord.Clear()
			nonce, err := blockchain.NewNonce()
			if err != nil {
				return err
			}
			dc.Bc.AddNewBlock(content, nonce)
		}
	}
	return nil
}
