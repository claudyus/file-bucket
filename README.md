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
| 403 | token doesn't exist |
| 412 | missed file |
| 405 | file yet exist |
| 401 | cannot create bucket dirs |

