# EZPwd Encrypted password manager compatible with GnuPGP/OpenPGP

A small utility that helps to keep passwords organized in an encrypted file.

The password file format:

```
service / username / password / comment 
```

Easy to grep. 
Also could be read by `gpg -d ${filename}` 

### Options:
```
ezpwd -h         
Usage of ezpwd:
  -add
        Add new password
  -passfile string
        Password file (default "private/test-pass.enc")
```

### Build

`go build app/ezpwd.go`
