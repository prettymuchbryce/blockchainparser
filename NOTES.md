Notes:
A few things.
1. We should be labeling outputs as their transaction types.
2. Ultimately transaction outputs can only be spent once, by one address. This means while examing an input, we could go back and check if the input is a PKSH, or a multisig transaction and retroactively set the public key depending on spender's address. (But how does this affect the plan to use go-routines to improve performance by reading one .dat file at a time ?)