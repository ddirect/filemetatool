# filemetatool
`filemetatool` is based on the `filemeta` library, which extracts metadata from a file stores it in an extended attribute.

The metadata includes the file hash, but also the file last modification timestamp; this is needed to detect if the file has been modified and thus the hash needs to be recalculated.

## Synopsis
`filemetatool -do OPERATION [OPTION]... PATH...`

## Supported operations
- **list:** lists the files with their hash, if available
- **refresh:** updates the file hashes if necessary
- **stat:** prints statistics about the files
- **scrub:** checks if the stored hash matches the actual hash of the files

## Example
```
filemetatool -do refresh <file_tree>
```
Calculates the file hashes in `file_tree` where necessary.
