# Lockgit

Storing secrets in git is dangerous and sometimes even considered a bad practice.
Yet, many people require a place to store secrets and git is a useful tool that we
are used to using. So - enter __Lockgit__, a tool to make it easy to use encryption
to safely store secrets in a git repository.

## Getting Started

Lockgit can be easily installed as a binary with either Homebrew on OSX or Linuxbrew
on Linux.

```
brew install jswidler/tap/lockgit
```

You can also easily build from source with `go get github.com/jswidler/lockgit`.

### List of commands

```
Usage:
  lockgit [command]

Available Commands:
  init        Initialize a lockgit vault
  set-key     Set the key for the current vault
  reveal-key  Reveal the lockgit key for the current repo
  delete-key  Delete the key for the current vault
  add         Add files to the vault
  rm          Remove files and globs from the vault
  status      Check if tracked files match the ones in the vault
  commit      Commit changes of tracked files to the vault
  open        Decrypt and restore secrets in the vault
  close       Delete plaintext secrets
  ls          List the files in the lockgit vault
  globs       List all saved glob patterns
  help        Help about any command
```


## Using Lockgit

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
 
Some of the files are too sensitive to check into git without encryption. Let's encrypt them with Lockgit
 

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
            FILE           | UPDATED |    PATTERN    |               HASH
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

##### Use source conntrol
You should check the entire `.lockgit` folder into source control.  

Lockgit can also update `.gitignore` as you use it, which helps prevent accidentally checking in your secrets.  `**/creds.json`, `**/*.pem` have both been added to it in our example

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
            FILE           | UPDATED |    PATTERN    |               HASH
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