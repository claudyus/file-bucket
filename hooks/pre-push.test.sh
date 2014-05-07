#!/bin/sh
#
# args:
# 1) bucket token
# 2) upload filename
# 3) file size
# 4) remote ip_client:port
#
echo $* > /tmp/pre-args.txt

#deny push on 3547264155b541de5cb6a7eae0431a14
if [[ $1 == '3547264155b541de5cb6a7eae0431a14' ]]; then
	exit 1;
fi

#allow overwrite on 9c7c64cfac3b922aacd4ba68cb631e82
if [[ $1 == '9c7c64cfac3b922aacd4ba68cb631e82' ]]; then
	exit 2;
fi

#allow all
exit 0
