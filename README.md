# clasp [![Build Status](https://travis-ci.org/vdemeester/clasp.svg?branch=master)](https://travis-ci.org/vdemeester/clasp)

`clasp` is a mini hook / rebuild configuration binary written in Go.

```bash
$ clasp $HOME/.gitconfig $HOME/.config/git/config.d
# […]
$ cat $HOME/.gitconfig
# Autogenerated from clasp at 2016-11-24T16:53:06+01:00


# Appending /home/vincent/.config/git/config.d/00-default
[alias]
    co = checkout
    st = status
    ci = commit --signoff
    cia = commit --amend
    ciad = commit --amend --date=\"$(date -R)\"
    ciads = commit --amend --date=\"$(date -R)\" -S
    civ = commit -v
```