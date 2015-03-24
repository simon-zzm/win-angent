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
# 上传批处理，并调用批处理删除上传文件。
files = {
        'file':('ttt.zip',open("E:\\del_file.bat", 'rb'),'application/octet-stream')
        }
_param = {'cline': "del_file.bat"}
res= requests.post("http://localhost:8866/comm/",  params = _param)
print res.url,res.text
