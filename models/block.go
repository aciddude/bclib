package models

import (
  "time"
)

type Hash256 []byte

type TxInput struct {
  Hash              Hash256     `db:"hash"`       // Hash of some previous transaction
  Index             uint32      `db:"index"`      // Output of the previous transaction
  Script            []byte      `db:"script"`     // Useless?
  Sequence          uint32      `db:"sequence"`   // Always 0xFFFFFFFF
  ScriptWitness     [][]byte
}

type TxOutput struct {
  Index             uint32      `db:"index"`      // Output index
  Value             int64       `db:"value"`      // Satoshis
  Hash160           []byte      `db:"hash160"`    // Public key
  Script            []byte      `db:"script"`     // Where the magic happens
}

type Transaction struct {
  NVersion          int32       `db:"n_version"`  // Always 1 or 2
  Hash              Hash256     `db:"tx_hash"`    // Transaction hash (computed)
  NVin              uint32      `db:"n_vin"`      // Number of inputs
  NVout             uint32      `db:"n_vout"`     // Number of outputs
  Vin               []TxInput
  Vout              []TxOutput
  Locktime          uint32      `db:"locktime"`
}

// BlockIndexRecord contains general index records parameters
// It defines the structure of the postgres table
type BlockHeader struct {
  NVersion          int32     `db:"n_version"`    // Version
  NHeight           uint32    `db:"n_height"`             // (Index)
  NStatus           uint32    `db:"n_status"`// (Index)
  NTx               uint32    `db:"n_tx"`// (Index)
  NFile             uint32    `db:"n_file"`// (Index)
  NDataPos          uint32    `db:"n_data_pos"`// (Index)
  NUndoPos          uint32    `db:"n_undo_pos"`// (Index)
  HashBlock         Hash256   `db:"hash_block"`// current block hash (Added)
  HashPrevBlock     Hash256   `db:"hash_prev_block"`// previous block hash (Index)
  HashMerkleRoot    Hash256   `db:"hash_merkle_root"`// (Index)
  NTime             time.Time `db:"n_time"`// (Index)
  NBits             uint32    `db:"n_bits"`// (Index)
  NNonce            uint32    `db:"n_nonce"`// (Index)
  TargetDifficulty  uint32    `db:"target_difficulty"`//
  Length            uint32
}

type Block struct {
  BlockHeader
  Transactions  []Transaction
}
