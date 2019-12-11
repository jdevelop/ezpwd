# EZPwd Encrypted password manager compatible with GnuPGP/OpenPGP

A small utility that helps to keep passwords organized in an encrypted file.

The password file format:

```
service / username / password / comment 
```

Easy to grep. 
Also could be read by `gpg -d ${filename}` 

More info on [https://ezpwd.jdevelop.com/](https://ezpwd.jdevelop.com/)

### Build & install

`go install ./...`
