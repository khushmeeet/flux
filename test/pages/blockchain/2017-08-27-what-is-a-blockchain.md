---
template: post
title: Intro to Blockchain
shortie: A simple guide on what is a blockchain and related cryptocurrencies
date: 2017-08-27
categories: cryptocurrency
tags:
  - blockchain
  - cryptography
---

Blockchain right now is the most hyped technology and is being touted as the
next step in the digitization of our economy, moving from centralized systems to
decentralized systems where no one controls the system. Blockchain is nothing
but a public ledger, where your transactions get recorded. First distributed
blockchain was conceived by Satoshi Nakamoto (no one knows who he is) in 2008
and in the following year implemented **Bitcoin** cryptocurrency with blockchain
serving as its public ledger. Since then bitcoin has grown exponentially and
current exchange for 1BTC (Bitcoin) is 4352.24 US Dollars. Many more
cryptocurrencies have also made their way into the market like **litecoin** or
**zcash** and many more. The second largest player in cryptocurrencies is
**ethereum**. Ethereum by the way is much more than just a blockchain and
cryptocurrency. It is a platform for building applications for that runs on
ethereum blockchain. Ofcourse we will discuss ethereum in a seperate article.

For now let's get down to how blockchain works and where do these currencies
come from.

### Securing the ledger
Now to understand how blockchain works, we need to understand our current system
and how blockchain solves the problems associated with it.

Right now we have a centralized entity, that records our transaction or we can
say acts a mediator between the sender and the receiver. The working of the
entity is opaque, we don't know how they operate, they may be corrupt (common
with government institutions) or may be biased to some customers.

Or OK!, let's forget about all this stuff and start with a very simple case of 4
people (Alice, Bob, John, Chris) trying to manage money splits among them.

They start by creating a public ledger, where everyone writes who need to pay to
whom and how much. Its an okay way to manage the debt, but this requires trust.
It's a public ledger which means that anybody can write to it. So maybe Bob
writes `John pays Bob 50$` without him telling. Now cryptography comes into
play. System of **Public Keys and Private Keys** are leveraged to produce and
authenticate **digital signature** of a person. Every person has a set of public/private key. Private key is kept secret and public key is well.. public. Whenever a new transaction is added to ledger, it gets signed by the person using the private key like this.
```
Sign(message, private key) -> Signature
```
These signatures are irreversible, meaning that you cannot derive message just
by looking at the signature. And such signatures are usually 256bit long (2^256 possibilities!!) hence
guessing is impossible. Therefore it is safe to say that these signatures can be
trusted.

Now to verify that the transaction is not forged, we use public key (anyone can
verify). Something like this
```
Verify(message, signature, public key) -> True/False
```
This verify functions tells us that the message is signed by the private key
associated with the public key that we are using.

Now to prevent a person from just copying the whole message again (because the
signature remains same for same message) a unique id is added with the message,
which becomes
```
Sign(message, unique id, private key) -> Signature
```
Now signature is changing with every transaction added to ledger, hence no
copying/forging possible.

But still, if you see our ledger system is centralized. This is where concept of
blockchain is introduced. The whole point of doing blockchain thing is to remove
trust on one single entity or to remove trust altogether and let the rules and
protocol do its work.

### Blockchain
As we know our ledger is just a table of recorded transactions. But to remove
that trust, to remove that centralization, everyone can hold a copy of
transaction. Hereby removing the need for handing over the ledger to someone for
maintenance.

To add a transaction, a broadcast will be made, so that every ledger remains in
sync (decentralized).

But we have a problem here, there is no guarantee that when a transaction is
made it will get added to every other ledger. Simply put there is no way to
authenticate that the broadcast was original. This is a very flawed way of
syncing the ledger, just listening to the broadcast.

This very problem was addressed in the original bitcoin paper. Here is a [link](https://bitcoin.org/bitcoin.pdf) to it.

General idea here is that the ledger with highest amount of computational word
is the true ledger. This computational work involves cryptographic hash
function. This means that for someone to produce a fraudulent transaction would
require an insane amount of computation, which is pretty infeasible.

The hash function generally used in cryptocurrencies is **SHA256**. SHA256 is
highly random in nature. Meaning, just a slight change in input results in
complete change of hash value. And even though there is no actual proof, SHA256
is considered irreversible. Input cannot be produced from the output.

Now to use SHA256 to define computational work, hash is calculated for the
ledger. The way it is done is that a condition is given on the hash value. Let's
say the condition here is that 30 zeros should be there at the starting of the hash. So we need to find a number that when added with the ledger and calculated the SHA256 of the whole thing gives us that hash with 30 zeros at the beginning. The only way this can be done is by iterating through every number and check whether its hash satisfies the condition. This is called the **proof of work** (PoW). Proof that this person has done work without the need to verify. This hash is tied to the ledger itself. Therefore any change in ledger, will change its hash and all the work needs to be done again.

Now to make a distributed tamper proof ledger, the ledger is divided into
blocks. Each blocks has list of transactions together with the proof of work.
The block is valid only when it has a proof of work. Now to maintain order and
timeliness to the blocks, hash of the previous block is added to the next block
hence SHA256 of a block becomes
```
SHA256(hash of previous block, transactions, PoW) -> Hash
```
Who calculates these hashes?. Miners do. Miners (block creators) collect the
transaction and then do the proof of work. When a new hash is created, meaning a
block is created, it gets added to the block chain. And the first person to do
so will get a payment of the that cryptocurrency. So with every block created,
the amount of currency will increase.

How this system is tamper proof?. Let's say Alice tries to fool Charlie by
adding a payment from Charlie to Alice. Therefore she hash to calculate the hash
of the block first to add it to the blockchain. But she has to do all of the
work herself to maintain that version of the blockchain. But will eventually be
overpowered by every other miner, as she would require more than 50% of the computing power to quickly calculate hash which isn't feasible.
Therefore her blockchain will fall short of the rest and therefore only the
longest blockchain prevails or is trustable.

So this is how a blockchain works. Blocks of transactions joined linearly in a
list. Everybody having his/her own copy of blockchain. Miners adding more
blocks every few minutes. Protected with a hash function.
This makes a pretty reliable decentralized system.

This is how bitcoin and several other currencies works. Though ethereum goes a
bit far by enabling blockchain for not just transactions, but any kind of
application.

A link to a fantastic video on blockchain and cryptocurrencies by
[3Blue1Brown](https://www.youtube.com/channel/UCYO_jab_esuFRV4b17AJtAw) -
[Ever wonder how Bitcoin (and other cryptocurrencies) actually work?](https://www.youtube.com/watch?v=bBC-nXj3Ng4).
This article is inspired by this video.
