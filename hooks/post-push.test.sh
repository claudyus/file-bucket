#!/bin/bash
#
# args:
# 1) bucket_token
# 2) complete_filename_with_path
# 3) remote ip_client:port
#

#test: echo info about uploaded file
echo $* > /tmp/post-args.txt

#print back to http
echo $*
