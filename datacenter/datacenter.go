package datacenter

import (
	"encoding/json"
	"errors"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/blockchain"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/client"
	"github.com/mcxxmc/simple-implementation-ecc/ecc"
	"github.com/mcxxmc/simple-implementation-ecc/galois"
	"strconv"
	"time"
)

type Record [512]byte	// assumes a record to be a byte array of 512 bytes.

// Msg a single message recording an activity
type Msg struct {
	Timestamp 		string	`json:"timestamp"`
	Entity			string	`json:"entity"`		// which is the user
	EntityPubKey	string	`json:"entity_pub_key"`
	Info 			string	`json:"info"`
}

// NewMsg returns a pointer to a new Msg object to be used by Record.
func NewMsg(entity int, entityPubKey *galois.Point, info string) *Msg {
	return &Msg{
		Timestamp: time.Now().String(),
		Entity: strconv.Itoa(entity),
		EntityPubKey: ecc.StringifyPublicKey(*entityPubKey),
		Info: info,
	}
}

// DataCenter the datacenter struct plays as the data center, which contains a blockchain holding past activity and
// records the public key of all the active users	//todo: incorporate a database to store both offline and online users
//
// Different from the design of bitcoin or other cryptocurrencies, users are not distinguished by there public key,
// but by an unique id; therefore, multiple users can possibly share the same public key.
//
// Please use NewDataCenter() as the constructor, and call RandomInitialization() to set a new private key and the public key.
// Then you can use RequestPublicKey() to get the corresponding public key.
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
func (dc *DataCenter) RandomInitialization() error {
	dc.Self.RandomInitialization()
	dc.Users[dc.Self.Id] = dc.Self.PubKey		// add own pub key to the users

	// write a msg
	err := dc.Write(NewMsg(dc.Self.Id, dc.Self.PubKey, InfoDatacenterReady))
	return err
}

// RequestPublicKey returns a copy of the public key given the user id.
func (dc *DataCenter) RequestPublicKey(userId int)	(galois.Point, error) {

	// currently, no msg for this

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

	// write a msg
	err := dc.Write(NewMsg(userId, &userPubKey, InfoRegisterNewUser))
	return err
}

// RequestUpdateUserPubKey updates the public key of a certain user; should verify the identity of the user first.
//
// sharedKeyFromUser = user privateKey * datacenter pubKey = datacenter privateKey * user pubKey
func (dc *DataCenter) RequestUpdateUserPubKey(userId int, sharedKeyFromUser, newUserPubKey galois.Point) error {
	userPubKey, exist := dc.Users[userId]
	if !exist {
		return errors.New("user id not found")
	}

	// write the first msg
	err := dc.Write(NewMsg(userId, userPubKey, InfoUpdateKey))
	if err != nil {
		return err
	}

	calculatedSharedKey := ecc.Calculate(*userPubKey, dc.Self.Ecdh.PrivateKey, dc.Self.Ecdh.Ep)
	if !galois.PointEqual(calculatedSharedKey, sharedKeyFromUser) {
		return errors.New("shared key does not match")
	}

	p := galois.Copy(newUserPubKey)
	dc.Users[userId] = &p

	// write the second msg
	err = dc.Write(NewMsg(userId, &newUserPubKey, InfoUpdateKeySuccess))

	return err
}

// Write writes bytes of data into the current record, and creates new block as needed.
func (dc *DataCenter) Write(data interface{}) error {		// todo: add encryption (e.g., using gcm)
	if data == nil {
		return errors.New("no data to write")
	}

	bytes, err := json.Marshal(&data)
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

// ForceWrite writes the current record into a new block, regardless if it is full.
func (dc *DataCenter) ForceWrite() error {	// todo: add encryption (e.g., using gcm)
	content := dc.ActiveRecord.Bytes()
	dc.ActiveRecord.Clear()
	nonce, err := blockchain.NewNonce()
	if err != nil {return err}
	dc.Bc.AddNewBlock(content, nonce)
	return nil
}
