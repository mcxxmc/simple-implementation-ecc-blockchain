# simple-implementation-ecc-blockchain
 A simple implementation for blockchain using my previous ECC algorithm. In Dev 2022.


It uses my previous package `github.com/mcxxmc/simple-implmentation-ecc` for the
calculation over finite field, the implementation of elliptic curve and the ECDH
process (including gcm encryption and decryption).


It consists of 3 parts: `blockchain`, `client` and `datacenter`.


`datacenter` is the simulation part. This project is not aimed at the 
simulation of any kind of cryptocurrency. Its main purpose is to simulate the 
usage of blockchain technology in the data storage, which is immutable and secure after 
creation by the nature of blockchain.


Users can register themselves at the "datacenter", and the "datacenter" is therefore
like a phone book. It distinguishes each user by their "phone number" (a unique id)
instead of their public keys (the public keys are used as unique ids by many 
cryptocurrencies, though).


TODO:

a complete verification system.

A real world example of identity verification using ecc is that,
a buyer from the other party can use your phone number (unique id) to verify your ownership 
of certain asset, without seeing the traditional official paper, which prevents 
potential leak of privacy. However, the buyer should also have registered in the
datacenter (owning a public key) and agrees on the same elliptic curve. The 
mathematics secret behind this mechanism is:


`buyer's privKey * your PubKey = buyer's privKey * your PubKey from datacenter = buyer's pubKey * your privKey`


note that during that exchange process, the buyer doesn't send the calculated shared key to the
seller; he/ she receives the shared key from the seller and compare. If they match, then the
seller's identity is verified. 


This is a very simple implementation and have fun!
