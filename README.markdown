cbdocloader
===========

cbdocloader is a small command line tool used to load a bunch of data
into a Couchbase cluster.

Build
-------

go build

Run command
------------

    cbdocloader OPTIONS DOCUMENTS

DOCUMENTS:

The documents parameter is a zip file containing all of the documents
to be uploaded.

The zip file must have the following layout:

    name/design_docs/  Which contains all the design documents for the views
    name/docs/         which contains all of the documents to be created

name may be some arbirary name and is completely ignored. The name of
the file in `/docs/` will be used for the key, and the entire body
will be used for the value. (if the filename have a .json suffix it
will be removed)


OPTIONS:

    -n HOST[:PORT]     Default port is 8091
    -u USERNAME        This parameter is ignored.
    -p PASSWORD        Bucket password
    -b BUCKETNAME      Specific bucket name. Default is default bucket. Bucket have to exist
    -s QUOTA           This parameter is ignored
    -h                 Show this help message and exit

Requirements
------------

The bucket needs to exist before trying to load the data.

Example
-------

Upload documents archived in zip file `../samples/gamesim.zip`. All data
will be inserted in bucket mybucket.


    ./cbdocloader -n localhost:8091 -u mybucket -p my -b mybucket ../samples/gamesim.zip
