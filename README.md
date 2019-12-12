There are a plethora of password managers, so why build another? I’ve got a few simple reasons why I think there is room for 1 more.
* Fancy UI is overrated; I just need something that works. Unix-way is the real thing.
* I want to use a single password manager across operating systems like Linux desktops, Linux/FreeBSD servers (and MacOS, since it is standard de-facto in many corporate environments ). And potentially Windows, but spare me this fate, please.
*  I want to have my password database as simple as possible. Essentially, just a good ol’ **text file**.
* I don’t want to risk getting locked out of my password manager if the app crashes or is not supported anymore. Hence I want my password database to be compatible with something that I can use pretty much anywhere. In fact, it must be so simple that I **shouldn't ever need a password manager app to access my passwords at all**.

So I have narrowed down my design considerations which consist of:
* text file
* encrypted by GnuPG
* with format of `service / username|email|anything / password / comment`. This format I’ve been using for more than 15 years myself, so it is time-proof. At least for me.
* [Golang](https://golang.org/) - clean, simple, better C, great language ( generics though... but who cares? )
* [Tview](https://github.com/rivo/tview) - the UI for the people who are not really comfortable with the hardcore command line and ASCII art.

An example ( unencrypted ) text file with the passwords would look like:
```
github / user@domain.com / ObHivyasvoHas0 / primary github account
atlassian / user+atlassian@gmail.com / Rud8Vor.Drivinn / JIRA Account etc
```

## Text interface

The basic user interface could be accessed by invoking `ezpwd` command in a terminal ( presuming that it is either in `$PATH` or in current directory ).
Help is available with `-h` command line option:
```
ezpwd -h
Usage of ezpwd:
  -add
        Add new password
  -list
        List all passwords
  -passfile string
        Password file (default "private/test-pass.enc")
  -update
        Update password
```

To start, the file needs to be created, and it could be achieved by `-add` switch. `ezpwd` will ask for the storage password ( this is the first point of failure - the password needs to be strong enough and should be memoized - not written down. It will be used to decrypt the storage with other passwords ). The storage file will be created in `$HOME/private/test-pass.enc` so make sure that this folder exists.
```
ezpwd -add
Storage Password :/>
Service :/> github
Username/email :/> user@domain.com
Enter Password :/>
Confirm Password :/>
Comment :/> primary github account

```
check out the content of the folder: 
```
ls -l ~/private/
-rw-r--r-- 1 user user 173 Dec 1 10:08 test-pass.enc
```
this is the encrypted file. What is cool about it is that you can easily decrypt it with `gpg`:
```
gpg -d ~/private/test-pass.enc
gpg: AES encrypted session key
gpg: encrypted with 1 passphrase
github / user@domain.com / ObHivyasvoHas0 / primary github account
```
Substantially, you don't need to worry about `ezpwd`. If you don't have it installed - you may use GnuPG to access the passwords. And AES is pretty good encryption protocol.

Now, if you want to use the password database - you just invoke `ezpwd` with no parameters:

```
ezpwd
Storage Password :/>
+---+---------+-----------------+------------------------+
| # | SERVICE |      LOGIN      |        COMMENT         |
+---+---------+-----------------+------------------------+
| 0 | github  | user@domain.com | primary github account |
+---+---------+-----------------+------------------------+
Choose password
```

Once you enter the storage password - `ezpwd` will decrypt the storage and show you a table. This table doesn't have the passwords listed. Instead, `ezpwd` will **copy the password into the clipboard** so it will be accessible from there and you can easily paste the password into the corresponding input field. This is not too much of security here, since any application may intercept the changes in the clipboard and steal passwords. But yet it is more convenient rather than select the password and copy it manually.

To choose the password - type the appropriate number, in this case `0` - as listed in the leftmost column. `ezpwd` will copy it into the clipboard and then exit. You can verify it by hitting `Ctrl-V` in the terminal window - that will insert the password for the service from the clipboard.

Update passwords works in the same way: you need to run `ezpwd -update`, then provide the **storage password** to decrypt the storage file, then specify the number corresponding to the password entry.
```
ezpwd -update
Storage Password :/>
+---+---------+-----------------+------------------------+
| # | SERVICE |      LOGIN      |        COMMENT         |
+---+---------+-----------------+------------------------+
| 0 | github  | user@domain.com | primary github account |
+---+---------+-----------------+------------------------+
Please choose the entry you'd like to change: 0
Service :/> github.com
Username/email :/> 
Enter Password :/> 
Confirm Password :/> 
Comment :/> 
```
Provide the password index (`0`) and then you may update some of the information. If you don't want to update certain field - just press `Enter` and this value will remain unchanged. In the example above the name of the service was changed, the rest of the fields weren't updated. This can be verified by 
```
ezpwd
Storage Password :/>
+---+------------+-----------------+------------------------+
| # |  SERVICE   |      LOGIN      |        COMMENT         |
+---+------------+-----------------+------------------------+
| 0 | github.com | user@domain.com | primary github account |
+---+------------+-----------------+------------------------+
```
or
```
gpg -d ~/private/test-pass.enc
gpg: AES encrypted session key
gpg: encrypted with 1 passphrase
github.com / user@domain.com / ObHivyasvoHas0 / primary github account
```
