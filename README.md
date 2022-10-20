# gallery-daemon
The daemon for the gallery app

## Usage
### Indexation only
```
$ ./gallery-daemon index --help
Starts the daemon in index mode only.

Indexes the given folder and create a database file.

Usage:
  gallery-daemon index [flags]

Flags:
  -h, --help   help for index

Global Flags:
  -f, --folder string   The folder containing the photos (default ".")
      --local-db        Place the database in the current folder
      --verbose         Verbose display
```

### Serve
```
$ ./gallery-daemon serve --help
Starts the daemon in server mode.

In this mode it will :
 - index the images if the database is not up-to-date
 - register the daemon to the backend
 - watch for file changes
 - serve backend requests

Usage:
  gallery-daemon serve [flags]

Flags:
  -H, --external-host string   External host (default "localhost")
  -h, --help                   help for serve
  -n, --name string            Daemon name (default "localhost-daemon")
  -p, --port int32             Grpc Port (default 9000)
      --re-index               Launch a full re-indexation

Global Flags:
  -f, --folder string   The folder containing the photos (default ".")
      --local-db        Place the database in the current folder
      --verbose         Verbose display
```

# Exemple
```
$ ./gallery-daemon serve -p 9001 -f ~/Images/Photos
✓ Re-indexing folder /home/spyder/Images/Photos
✓ Done.
✓ Listening on 0.0.0.0:9001
✓ Daemon registered.
✓ Watching folder /home/spyder/Images/Photos
```
