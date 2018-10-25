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
  reveal-key  Reveal the lockgit key for the current repo
  unlock      Set the key for the current vault
  lock        Delete the key for the current vault
  ls          List the files in the lockgit vault
  add         Add files to the vault
  commit      Commit changes of tracked files to the vault
  rm          Remove files from the vault
  open        Decrypt and restore secrets in the vault
  close       Delete plaintext secrets
  status      Check if the secrets present match the ones in the vault
  help        Help about any command
```


## Using Lockgit

Suppose there is a small project with the following files in it.  

```
myserverconfig
├── config
│   ├── tls.ca
│   ├── tls.crt
│   └── tls.key
├── creds
│   └── awscreds
├── scripts
│   ├── init.sh
│   └── run.sh
└── templates
    └── nginx.conf.tpl
 ````
 
 Two of the files have secrets in them, `config/tls.key` and `creds/awscreds`.  Let's put
 them in a Lockgit vault.  
 

##### Initialize a vault
First, initialize a new vault in the `myserverconfig` directory:
 
```
$ lockgit init
Initialized empty lockgit vault 'flattered-gusto' in /home/myserverconfig/.lockgit
```

If there are no parameters to init, lockgit will make up a name to refer to the repo.  This can be changed later with `rename`.

##### Add secrets
Next, add the secrets to it

```
$ lockgit add creds/awscreds config/tls.key
added creds/awscreds to vault
added config/tls.key to vault
```
We can see what secrets are in the vault with either `lockgit ls` or `lockgit status` .

The files have been encrypted and stored in the `.lockgit` directory.  It currently looks
something like this:

```
.lockgit/
├── data
│   ├── uJl28qpmje-cWGDUv3p8iiJgcANTPKdK
│   └── yocrUblPqDoRaKYnjE6VfQLQo6LrlrHT
├── lgconfig
└── manifest
``` 

##### Use source conntrol
You should check the entire `.lockgit` folder into source control.  

Lockgit can also update `.gitignore` as you use it.  We now have two secrets in
the `myserverconfig` folder: `config/tls.key`, and `creds/awscred`.
Both of these files have been added to `.gitignore` to help prevent them from being
checked in.

##### Delete and Restore plaintext secrets
Delete and restore your secrets with `lockgit close` and `lockgit open`.

##### Share the key with someone else
To see the key, use 
```
$ lockgit reveal-key
FA633KF422AXETBBMXUZYNXZDXN4VRKSE4TI4N2KTXYHV6MUAHQA
```  

To use this key to unlock the vault, use `unlock`

```
$ lockgit unlock FA633KF422AXETBBMXUZYNXZDXN4VRKSE4TI4N2KTXYHV6MUAHQA
```

The key is saved to your home directory in the config file `~/.lockgit.yml` (unless you overrode this location from
the command line).  You can remove the key from the config file by using `lock`.  Be wary that this will delete your key,
so if it isn't written down somewhere, you will lose the contents of the vault.


##### Make changes to your secrets
After you update a secret, lockgit can detect the change.
```
$ lockgit status
       FILE      | UPDATED |               HASH
+----------------+---------+----------------------------------+
  config/tls.key | true    | uJl28qpmje-cWGDUv3p8iiJgcANTPKdK
  creds/awscreds | false   | yocrUblPqDoRaKYnjE6VfQLQo6LrlrHT
```

To update the encrypted secret, first use `lockgit commit`
```
$ lockgit commit
/home/myserverconfig/config/tls.key updated
```

Then commit the changes to source control.  In this case there will be three changes:
```
deleted:    .lockgit/data/uJl28qpmje-cWGDUv3p8iiJgcANTPKdK
new file:   .lockgit/data/ZnrcquH2lU9KF6HW0E8oi4PGEISvqvc3
modified:   .lockgit/manifest
```

The two files in `.lockgit/data` are the encrypted secrets.

The manifest is a text file that can be easily examined.  This makes it possible to
see what secrets people are changing when reviewing commits.

```
$ cat .lockgit/manifest
ZnrcquH2lU9KF6HW0E8oi4PGEISvqvc3	config/tls.key
yocrUblPqDoRaKYnjE6VfQLQo6LrlrHT	creds/awscreds
```