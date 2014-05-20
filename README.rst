File Bucket
===========

An http server, with write only permission for backup solution.

The idea is to allow clients to push backup files to centralized server over LAN, unlike other backup solution file-bucket doesn't use the classic username/password pair to auth the user but a single 32-byte long token.

To enhanced security each server can **only** write into it's own bucket, so if a server is compromised the attacker can only guess the server token but it cannot read old backups neither discover or read other server backup files.

Installation
------------

To install the file-bucket on debian system you can dowload the package from http://claudyus.github.io/file-bucket/ install it with ``dpkg``::

  $ wget http://claudyus.github.io/file-bucket/file-bucket_VERSION.deb
  $ dpkg -i file-bucket_VERSION.deb

Configuration
^^^^^^^^^^^^^

You should configure the buckets list and buckets home inside ``/etc/file-bucket/config.json``
To override this settings you can pass ``--config <file>`` as command line parameters.

The bucket list can also be reload using ``SIGUSR2``, please note that only changes to bucket list are effective using this method.

Use
^^^

On your client you can simple push backup file using curl like:
::

  $ curl -X POST <serverip>:1234/<bucket_token> -F file=@<file>

The server will responde with ``200`` on success or with the following errorcode are used:

+-----------+------------------------------------------+
| ErrorCode | Detail                                   |
+===========+==========================================+
| 403       | given token doesn't exist                |
+-----------+------------------------------------------+
| 412       | missed file in form post                 |
+-----------+------------------------------------------+
| 405       | file yet exist on filesystem             |
+-----------+------------------------------------------+
| 401       | cannot create bucket dirs                |
+-----------+------------------------------------------+
| 409       | abort due to pre-push script (see below) |
+-----------+------------------------------------------+

Extending the behaviour
-----------------------

``file-bucket`` support a series of hooks that can be used to extend the default behaviour.
For example you can email someone when a file is pushed or deploy the file as a lxc container or the contained files in an automatic way.

At the moment the following hooks are defined:

- pre-push.sh
- post-push.sh

The ``/etc/file-bucket/pre-push.sh`` script is executed before the check of the file presence on the filesystem is executed.
This script should return **0** to allow the file upload, should return **1** to abort the upload and could return **2** to allow file overwrite.
This script is called with the arguments:

- bucket token
- upload filename
- file size (TODO)
- remote ip_client:port

The ``/etc/file-bucket/post-push.sh`` script is executed after the file upload and can be used to manipulate the file itself.
This script is called with the arguments:

- bucket token
- complete filename with path
- remote ip_client:port

The stdout of this command is passed back to the client as http body. Http return code is than set to 200.

**Note: All hooks script should have executable flag set.**

Hooks examples
^^^^^^^^^^^^^^

Some scripts files as example are availables inside the ``hooks/`` directory.

LICENSE
-------

The MIT License (MIT)

Copyright (c) 2014 Claudio Mignanti <c.mignanti@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
