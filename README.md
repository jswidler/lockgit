# Lockgit

Storing secrets in git is dangerous and sometimes even considered a bad practice.
Yet, many people require a place to store secrets and git is a useful tool that we
are used to using. So - enter __Lockgit__, a tool to make it easy to use encryption
to safely store secrets in a git repository.


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
Initialized empty lockgit vault in /home/myserverconfig/.lockgit
```

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
├── key
└── manifest
``` 

##### Use source conntrol
You should check the `.lockgit` folder into source control, with the exception of the `key` file.  

Lockgit can also update `.gitignore` as you use it.  We now have three secrets in
the `myserverconfig` folder: `.lockgit/key`, `config/tls.key`, and `creds/awscred`.
All of these files have been added to `.gitignore` to help prevent them from being
checked in.

##### Delete and Restore plaintext secrets
Delete and restore your secrets with `lockgit close` and `lockgit open`.

##### Share the key with someone else
To see the key, use 
```
$ lockgit reveal-key
FA633KF422AXETBBMXUZYNXZDXN4VRKSE4TI4N2KTXYHV6MUAHQA
```  

To use this key to recreate the `key` file, use `unlock`

```
$ lockgit unlock FA633KF422AXETBBMXUZYNXZDXN4VRKSE4TI4N2KTXYHV6MUAHQA
```

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

The manifest is a text file can be easily examined.  This makes it possible to understand what
secrets people are changing.

```
$ cat .lockgit/manifest
ZnrcquH2lU9KF6HW0E8oi4PGEISvqvc3	config/tls.key
yocrUblPqDoRaKYnjE6VfQLQo6LrlrHT	creds/awscreds
```