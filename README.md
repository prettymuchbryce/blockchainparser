>Parse the bitcoin blockchain into a postgresql database.

##Work in progress
This tool is not complete, and is still a work in progress.

####Description 
This program uses the Bitcoin .dat file format to parse the blockchain into a postgresql database for easier querying. This program parses Transactions, and Blocks into four postgresql tables called transactions, blocks, inputs, and outputs.

This program parses binary files on your harddrive. It won't auto-update, or be aware of new incoming transactions. You will have to re-run the program once those blocks have been added to your .dat files. In other words, this is for historical (non-realtime) data only.

####.DAT files
You will need to provide an argument with a link to your directory which contains these files. They look like this: blk00000.dat, blk00001.dat, etc.

If you're using Bitcoin QT, the .dat files can be located in these locations:

```
Linux:

~/.bitcoin/

MacOS:

~/Library/Application Support/Bitcoin/

Windows:

%APPDATA%\Bitcoin
```
####Setup

1. Install postgresql
2. Run schema.sh to setup the db and schema
3. Edit the configuration to specify the location of your .dat files
4. Compile and run main.go

####TODO
1. Parse P2SH transactions
2. Parse Multisig transactions
3. Create wallets table(?)

####Contributing
Pull requests are welcome and encouraged

#####Arguments
```
-data The path to your .dat file (blocks) directory
-user Your postgresql database user
```

####References
* [How to parse the bitcoin blockchain](http://codesuppository.blogspot.com/2014/01/how-to-parse-bitcoin-blockchain.html)
* [Bitcoin wiki](http://bitcoin.it)
* [A survey of transaction types](http://www.quantabytes.com/articles/a-survey-of-bitcoin-transaction-types)
* [Blockchain.info](http://www.blockchain.info)