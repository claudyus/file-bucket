# File Bucket

An http server, with write only permission for backup solution.

The idea is to allow clients to push backup files to centralized server over LAN, unlike other backup solution file-bucket doesn't use the classic username/password pair to auth the user but a single 32-byte long token.

To enhanced security each server can **only** write into it's own bucket, so if a server is compromised the attacker can only guess the server token but it cannot read old backups neither discover or read other server backup files.

## Installation

TODO

## Configuration
You should configure the buckets list and buckets home inside ```/etc/file-bucket/config.json```

## Use
On your client you can simple push backup file using curl like:

```
$ curl -X POST <serverip>:1234/<bucket_token> -F file=@<file>
```

The server will responde with ```200``` on success or the following errorcode are used:

| ErrorCode | Detail       |
| --------- | -------------|
| 403 | given token doesn't exist |
| 412 | missed file in form post |
| 405 | file yet exist on filesystem |
| 401 | cannot create bucket dirs |
| 409 | abort due to pre-push script (see below) |

## Extending the behaviour

```file-bucket``` support a series of hooks that can be used to extend the default behaviour.
For example you can email someone when a file is pushed or deploy the file as a lxc container or the contained files in an automatic way.

At the moment the following hooks are defined:
 * pre-push.sh
 * post-push.sh

The ```/etc/file-bucket/pre-push.sh``` script is executed before the check of the file presence on the filesystem is executed.
This script should return **0** to allow the file upload, should return **1** to abort the upload and could return **2** to allow file overwrite.
This script is called with the arguments:
 - bucket token
 - upload filename
 - file size (TODO)
 - remote ip_client:port

The ```/etc/file-bucket/post-push.sh``` script is executed after the file upload and can be used to manipulate the file itself.
This script is called with the arguments:
 - bucket_token
 - complete_filename_with_path
 - remote ip_client:port

The stdout of this command il passed back to the client as http body. Http return code is set to 200

**Note: All hooks script should have executable flag set.**

### Hooks examples
Some scripts files as example are availables inside the hooks/ directory.
