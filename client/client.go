package client

import (
	"github.com/mcxxmc/simple-implementation-ecc/ecc"
	"github.com/mcxxmc/simple-implementation-ecc/galois"
)

// Client the client instance for simulation
//
// Please use NewClient() as constructor, and call RandomInitialization() to set a new private key and the public key.
// Then you can use GetPublicKey() to get the corresponding public key.
type Client struct {
	Ecdh		*ecc.InstanceECDH		// the private key is included in this ECDH instance
	PubKey		*galois.Point			// the public key
	Name		string					// username (does not have to be unique)
	Id 			int						// the unique id
}

// NewClient returns a pointer to a new Client object.
//
// The elliptic curve must have been initialized with a generator point. (by ep.SetGeneratorPoint())
func NewClient(ep *ecc.Elliptic, name string, id int) *Client {
	ecdh := ecc.NewInstanceECDH(ep)
	return &Client{
		Ecdh: ecdh,
		Name: name,
		Id:   id,
	}
}

// RandomInitialization randomly sets a new private key and the corresponding public key.
//
// This can be used to update the keys as well.
func (client *Client) RandomInitialization() {
	client.Ecdh.RandomlyPicksPrivateKey()
	key := client.Ecdh.PublicKey()
	client.PubKey = &key
}

// GetPublicKey returns a copy of the public key.
//
// Can only be called after RandomInitialization() is called, and you should always call this method for getting public key.
func (client *Client) GetPublicKey() galois.Point {
	return galois.NewPoint(client.PubKey.X, client.PubKey.Y)
}

// EncryptMsg encrypts the msg
func (client *Client) EncryptMsg(msg string) (*ecc.EncryptedMsg, error) {
	return client.Ecdh.Encrypt(msg, client.GetPublicKey())
}

// DecryptMsg decrypts the msg
func (client *Client) DecryptMsg(encrypted *ecc.EncryptedMsg) (string, error) {
	return client.Ecdh.Decrypt(encrypted)
}
