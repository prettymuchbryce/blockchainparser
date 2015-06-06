For the bold and the brave, a golang program which parses the bitcoin blockchain into a postgresql database.

At the moment this program only parses Transactions, Blocks, and Wallets, however; it could be easily modified to do whatever you want and pull requests are welcome.

This program uses the Bitcoin QT .dat file format to parse the blockchain. You will need to provide an argument with a link to your directory which contains these files. They look like this: blk00000.dat, blk00001.dat, etc.

This program parses binary files on your harddrive. It won't auto-update, or be aware of new incoming transactions. You will have to re-run the program once those blocks have been added to your .dat files. In other words, this is for historical (non-realtime) data only.

If you're using Bitcoin QT, the .dat files can be located in these locations:

Linux:

~/.bitcoin/

MacOS:

~/Library/Application Support/Bitcoin/

Windows:

%APPDATA%\Bitcoin

Setup

1. Install postgresql
2. run schema.sh to setup the db and schema
3. Edit the configuration to specify the location of your .dat files
4. compile and run main.go

Arguments

-data The path to your .dat file directory
-limit A block # limit for the parsing. Useful if you just want to try parsing the first n blocks to see how the data looks

Special Thanks to John Ratcliff for his exemplary in-depth writeup of the blockchain binary format.

References
How to parse the bitcoin blockchain - http://codesuppository.blogspot.com/2014/01/how-to-parse-bitcoin-blockchain.html