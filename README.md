# ddSync Agent

ddSync stands for Distributed Directory Synchronization.

It was originally written to keep a number of directories in sync across a few different machines.

## How it works
 
- A simple syncing agent that monitors.
- The application communicates using nsq.
- Events on the file system are enriched with more information about the event.
- Events are published on an nsq topic.
- Other agents 

## How to deploy

### Running from the binary

- Run an nsq lookup instance.
- Run an nsq admin instance (if you choose).
- Run ddsync with `-dir` flag indicating directories comma separated.
- `ddsync -dirs=.` for the local directory.
- `ddsync -dirs=~/sync` for the sync directory in the user's home folder.
- `ddsync -dirs=~/sync,~/other` for the sync and other directories in the user's home folder.


### Running docker container

TBA
