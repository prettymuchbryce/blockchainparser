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

psql -d blockchainparser -c "CREATE TABLE transactions (
    hash varchar(64) primary key,
    version bigint NOT NULL,
    lock bigint NOT NULL
);"

psql -d blockchainparser -c "CREATE TABLE inputs (
    hash varchar(64) references transactions(hash),
    index bigint NOT NULL,
    script bytea NOT NULL,
    sequence bigint NOT NULL
);"

psql -d blockchainparser -c "CREATE TABLE outputs (
    hash varchar(64) references transactions(hash),
    value varchar(64) NOT NULL,
    script bytea NOT NULL
);"