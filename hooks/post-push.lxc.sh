#!/bin/bash
#
# args:
# 1) bucket_token
# 2) complete_filename_with_path
# 3) remote ip_client:port
#

# concept/untested
if [[ $1 == '9c7c64cfac3b922aacd4ba68cb631e82' ]]; then
	sudo lxc-stop --name my-lxc
	sudo tar xvzf $2 -C /var/lib/lxc/
	sudo lxc-start --name my-lxc
	#rm $2
fi
