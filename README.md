<img src="./images/lockgit.png" alt="LockGit" width="400px">

[![License](https://img.shields.io/github/license/jswidler/lockgit)](https://github.com/jswidler/lockgit/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/jswidler/lockgit)](https://github.com/jswidler/lockgit/releases)
[![CircleCI](https://img.shields.io/circleci/build/github/jswidler/lockgit/master)](https://circleci.com/gh/jswidler/lockgit/tree/master)
[![Go Report](https://goreportcard.com/badge/github.com/jswidler/lockgit)](https://goreportcard.com/report/github.com/jswidler/lockgit)
[![Go Lang](https://img.shields.io/github/go-mod/go-version/jswidler/lockgit)](https://golang.org)

LockGit is intended to be used with source control to safely store secrets using encryption,
so the secrets are unreadable to anyone who has access to the repository but does not have the key.

#

## Table of Contents

* [Getting Started](#getting-started)
  * [Installation](#installation)
  * [List of commands](#list-of-commands)
* [Using LockGit](#using-lockgit)
  * [Initialize a vault](#initialize-a-vault)
  * [Add secrets](#add-secrets)
  * [Use source control](#use-source-control)
  * [Delete and Restore plaintext secrets](#delete-and-restore-plaintext-secrets)
  * [Share the key with someone else](#share-the-key-with-someone-else)
  * [Make changes to your secrets](#make-changes-to-your-secrets)
* [Security](#security)
  * [Encryption](#encryption)
  * [Files](#files)
  * [Other safety](#other-safety)

## Getting Started
Storing secrets in Git is dangerous and sometimes even considered a bad practice.
Yet, many people require a place to store secrets and git is a useful tool that we
are used to using. So - enter __LockGit__, a tool to make it easy to use encryption
to safely store secrets in a Git repository.

### Installation

LockGit can be installed as a binary with either Homebrew on OSX or Linuxbrew on Linux.

```
brew install jswidler/tap/lockgit
```

You can also build from source with `go get github.com/jswidler/lockgit`.  Bash and zsh completion are installed for you if you use brew, so that is the preferred method.

### List of commands

```
Usage:
  lockgit [command]

Available Commands:
  init        Initialize a lockgit vault
  set-key     Set the key for the current vault
  reveal-key  Reveal the lockgit key for the current repo
  delete-key  Delete the key for the current vault
  add         Add files and glob patterns to the vault
  rm          Remove files and globs patterns from the vault
  status      Check if tracked files match the ones in the vault
  commit      Commit changes of tracked files to the vault
  open        Decrypt and restore secrets in the vault
  close       Delete plaintext secrets
  ls          List the files in the lockgit vault
  globs       List the saved glob patterns in the vault
  help        Help about any command
```


## Using LockGit

Suppose there is a small project with the following files in it.  

```
myserverconfig
├── config
│   ├── config.yml
│   ├── creds.json
│   └── tls
│       ├── cert.pem
│       ├── chain.pem
│       ├── fullchain.pem
│       └── privkey.pem
└── nginx.conf
```
 
Some of the files are too sensitive to check into Git without encryption. Let's encrypt them with LockGit.
 

##### Initialize a vault
First, initialize a new vault in the `myserverconfig` directory:
 
```
$ lockgit init
Initialized empty lockgit vault in /home/myserverconfig/.lockgit
Key added to /Users/jesse/.lockgit.yml
```

##### Add secrets
Next, add the secrets to it

```
$ lockgit add '**/creds.json' '**/*.pem'
added file 'config/creds.json' to vault
added file 'config/tls/chain.pem' to vault
added file 'config/tls/cert.pem' to vault
added file 'config/tls/privkey.pem' to vault
added file 'config/tls/fullchain.pem' to vault
added glob pattern '**/*.pem' to vault
added glob pattern '**/creds.json' to vault
```
We can see what secrets are in the vault with either `lockgit ls` or `lockgit status` .

```
$ lockgit status
            FILE           | UPDATED |    PATTERN    |                ID
+--------------------------+---------+---------------+----------------------------------+
  config/creds.json        | false   | **/creds.json | Oov8Rpf2YOU0mEQhGlHeDCzFHXRtkFnu
  config/tls/cert.pem      | false   | **/*.pem      | miehMYgqYtIVGMpVnss4ZZzlAQRpZAVd
  config/tls/chain.pem     | false   | **/*.pem      | m4_U5mtAOlEuXL5raxvWHRxBq2vq24Q3
  config/tls/fullchain.pem | false   | **/*.pem      | a1r4uoyv0XQpeltE7NjWD_93ufb27gzK
  config/tls/privkey.pem   | false   | **/*.pem      | BT19Sb8kQxx5Ztp20cX4IJQEAJE5vAkp
```


The files have been encrypted and stored in the `.lockgit` directory.  It currently looks
something like this:

```
.lockgit/
├── data
│   ├── BT19Sb8kQxx5Ztp20cX4IJQEAJE5vAkp
│   ├── Oov8Rpf2YOU0mEQhGlHeDCzFHXRtkFnu
│   ├── a1r4uoyv0XQpeltE7NjWD_93ufb27gzK
│   ├── m4_U5mtAOlEuXL5raxvWHRxBq2vq24Q3
│   └── miehMYgqYtIVGMpVnss4ZZzlAQRpZAVd
├── lgconfig
└── manifest
``` 

##### Use source control
You should check the entire `.lockgit` folder into source control.  

LockGit can also update `.gitignore` as you use it, which helps prevent accidentally checking in your secrets.  `**/creds.json`, `**/*.pem` have both been added to it in our example

##### Delete and Restore plaintext secrets
Delete and restore your secrets with `lockgit close` and `lockgit open`.

##### Share the key with someone else
To see the key, use 
```
$ lockgit reveal-key
FA633KF422AXETBBMXUZYNXZDXN4VRKSE4TI4N2KTXYHV6MUAHQA
```  

To use this key to unlock the vault, use `set-key`

```
$ lockgit set-key FA633KF422AXETBBMXUZYNXZDXN4VRKSE4TI4N2KTXYHV6MUAHQA
```

The key is saved to your home directory in the config file `~/.lockgit.yml` (unless you
overrode this location from the command line).  You can remove the key from the config
file by using `delete-key`.  Be wary that this will delete your key, so if it isn't written
down somewhere, you will lose the contents of the vault.


##### Make changes to your secrets
After you update a secret, lockgit can detect the change.

```
$ lockgit status
            FILE           | UPDATED |    PATTERN    |                ID
+--------------------------+---------+---------------+----------------------------------+
  config/creds.json        | true    | **/creds.json | 2HDEn74HAAws-D1Y2HS1ak7e0xGSo7kN
  config/tls/cert.pem      | false   | **/*.pem      | miehMYgqYtIVGMpVnss4ZZzlAQRpZAVd
  config/tls/chain.pem     | false   | **/*.pem      | m4_U5mtAOlEuXL5raxvWHRxBq2vq24Q3
  config/tls/fullchain.pem | false   | **/*.pem      | a1r4uoyv0XQpeltE7NjWD_93ufb27gzK
  config/tls/privkey.pem   | false   | **/*.pem      | BT19Sb8kQxx5Ztp20cX4IJQEAJE5vAkp
```

To update the encrypted secret, first use `lockgit commit`

```
$ lockgit commit
config/creds.json updated
```

Then commit the changes to source control.  In this case there will be three changes:

```
deleted:    .lockgit/data/2HDEn74HAAws-D1Y2HS1ak7e0xGSo7kN
new file:   .lockgit/data/Ik0gMeLDyIsIZNmNIEoeLzuH22kG2Cdp
modified:   .lockgit/manifest
```

The two files in `.lockgit/data` are the encrypted secrets.

The manifest is a text file that can be easily examined.  This makes it possible to
see what secrets people are changing when reviewing commits.

```
$ cat .lockgit/manifest
Ik0gMeLDyIsIZNmNIEoeLzuH22kG2Cdp	config/creds.json
miehMYgqYtIVGMpVnss4ZZzlAQRpZAVd	config/tls/cert.pem
m4_U5mtAOlEuXL5raxvWHRxBq2vq24Q3	config/tls/chain.pem
a1r4uoyv0XQpeltE7NjWD_93ufb27gzK	config/tls/fullchain.pem
BT19Sb8kQxx5Ztp20cX4IJQEAJE5vAkp	config/tls/privkey.pem
```

## Security

### Encryption

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

### Files

Most of the data LockGit will access on your filesystem will be inside of the project root, which is the location
where you initialize a LockGit directory. Generally this would also be the same root directory as the Git repository.
Inside the project root folder, LockGit will create a folder called `.lockgit`, which is intended to be checked into
source control. All the data in this folder is either not sensitive or encrypted.

The file outside the project root that LockGit will use is a file called `.lockgit.yml` which will be placed into your home directory (`~`).
The keys to each vault will be stored in this file.  If you read the YAML file, you will see a key and a path for each
vault.  The path is not important - it is only there to make it easier to identify the vault for a human.  The vault is actually identified by the UUID and the path in `.lockgit.yml` will update to the last known location of the vault.

### Other safety

The following points are provided to give assurance LockGit will never send data and that future updates will be 
backwards compatible.

- All official public releases will always be able to read secrets saved by older versions, so there is never a danger
you will lose access to your secrets by updating.
- LockGit has no network functionality.  It does not collect usage statistics, crash reports, your browser history, or
anything else, so it cannot leak your key and respects your privacy.
- LockGit will always be free and the source code is available under MIT license.