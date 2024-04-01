# Fzfnote
![Example](https://s3.ezgif.com/tmp/ezgif-3-d11979b696.gif)

Small terminal utility to add, read and delete one line notes. Inspired by Cleber's [pilha.sh](https://github.com/cleberzavadniak/pilha.sh)

## Reason:
Add, read and delete notes stored in a text file from the terminal. This program uses fzf to list notes and fuzzy find them. 

## Installation

```bash
# Make sure you have go installed
go install github.com/frnsimoes/fzfnote-go@latest
```

By default, `notes.md` is stored in users home directory `~/notes.md`. 


```bash

## Usage

- `fzfnote add text argument` - Add notes. One liners. Example: `fzfnote add This is a note`.
- `fzfnote read` - Read notes. Select notes will be copied to the clipboard. 
- `fzfnote delete` - Delete notes. `delete` will pipe the file to fzf and you can select the note to delete. You can select more than one note to delte by using fzf's <Tab>.
```

## Roadmap (Wishlist)

- [ ] Add / read / delete notes by tag
- [ ] Multiple file notes support
- [ ] Add / read / delete notes by creation date
