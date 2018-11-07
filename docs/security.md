<img src="../images/lockgit.png" alt="LockGit" width="400px">

# Security Overview

LockGit is intended to be used in conjunction with Git to safely store secrets in source control in
a way that makes those secrets unreadable to anyone who has access to the Git repository, but does not have the key.

Most of the data LockGit will access on your filesystem will be inside of the project root, which is the location
where you initialize a LockGit directory. Generally this would also be the same root directory as the Git repository.
Inside the project root folder, LockGit will create a folder called `.lockgit`, which is intended to be checked into
source control. All the data in this folder is either not sensitive or encrypted.

The other location LockGit will use is a file called `.lockgit.yml` which will be placed into your home directory (`~`).
The keys to each vault will be stored in this file.  If you read the YAML file, you will see a key and a path for each
vault.  The path is not important - it is only there to make it easier to identify which vault is which.

## Encryption

LockGit works using by saving data files in the `.lockgit/data` directory with 256 bit AES encryption in CFB mode. Each
encrypted file contains the contents of one file in the vault. The encrypted file also contains metadata with the
relative path and permissions of the file which are used when recreating the file. The contents of the data files are
compressed with zlib before encrypted.

The AES initialization vector is randomized each time a file is encrypted; therefore a different file is produced each
time a file is encrypted even if the contents are the same. Because the relative path is also stored in the encrypted
file, these files cannot be reused if a file moves, but is not changed. This is by design; so that edits to the
manifest cannot cause the secrets to end up in unexpected places.

A key to a LockGit vault is a 256 bit AES key. In text form, it is a 52 character base32 encoded string.

2<sup>256</sup> (about 10<sup>77</sup>) key possibilities is a lot. There are about 2<sup>80</sup> (10<sup>21</sup>)
stars in the observable universe - so 2^256 is, like, a really big number. AES is considered secure and uncrackable.  No
one will be able to decrypt the files without the key.

## Other safety

The following points are provided to give assurance LockGit will never send data and that future updates will be 
backwards compatible.

- All official public releases will always be able to read secrets saved by older versions, so there is never a danger
you will lose access to your secrets by updating.
- LockGit has no network functionality.  It does not collect usage statistics, crash reports, your browser history, or
anything else, so it cannot leak your key and respects your privacy.
- LockGit will always be free and the source code is available under MIT license.
