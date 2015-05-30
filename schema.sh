psql -c "DROP DATABASE blockchainparser"

createdb blockchainparser

psql -d blockchainparser -c "CREATE TABLE blocks ( 
    hash varchar(64) primary key,
    length bigint NOT NULL,
    version bigint NOT NULL,
    previousBlockHash varchar(64) NOT NULL,
    merkleRoot varchar(64) NOT NULL,
    timestamp bigint NOT NULL,
    difficulty bigint NOT NULL,
    nonce bigint NOT NULL,
    transactionCount numeric(20) NOT NULL
);"