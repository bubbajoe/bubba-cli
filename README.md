# bubba-cli (WIP)

## Searching

### `search`

searches for pattern in specific file(s)

#### USAGE

```bash
bb search $TARGET
bb search -d $PATH $TARGET # default $PATH="."
bb search -r $TARGET # r () search through all dire
bb search -R ^.*\.txt$ $TARGET # -r means TARGET must be a valid golang regular expression
bb search -n ^.*\.txt$ $TARGET # n (name regex) search through all files name
```

### `vsm`

Create an index which used VSM and inteactively search through that index.

## Encode/Decode

### `b64`

encodes a string to base64

### `b64d`

decodes a string from base64

### `urle`

encodes a string to percent-encoding

### `urld`

decodes a string from percent-encoding

## Networking

### `dns`

dns lookup

### `rdns`

reverse dns loopup


