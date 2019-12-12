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

For basic prompt-based text interface refer to [Text Interface](textview.md)
