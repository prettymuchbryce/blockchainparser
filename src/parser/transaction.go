package parser

//---
//Transaction version # (4 bytes)
//# of inputs (variable)
//  Hash
//  Transaction Index (4 bytes) uint32
//  Script length (variable)
//  Script data (length bytes)
//  Sequence # (4 bytes) uint32 always 0xFFFFFFFF
//Output count (variable)
//  Value (8 bytes) uint64
//  Script length (variable)
//  Output script (length bytes)
//lock time (4 bytes) uint32 always 0
//---

type Transaction struct {
	version uint32
	inputs  []*Input
	outputs []*Output
	lock    uint32
	hash    [32]byte
}

type Input struct {
	index    uint32
	script   []byte
	sequence uint32
}

type Output struct {
	value  uint64
	script []byte
}
