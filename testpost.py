#!/bin/env python
# -*- coding:utf-8 -*-
# Filename:
# Revision:
# Date:        2015-03-04
# Author:      王亮
### END INIT INFO

# 需要安装requests，在windows下执行
import requests
import json

files = {
        'file':('ttt.zip',open("E:\\ttt.zip", 'rb'),'application/octet-stream')
        }
data = {'md5':'xxxxxxxxxxxx'} //md5验证没有启动

res= requests.post("http://localhost:8866/upfile/",  files = files, params=data)
print res.url,res.text
