# !/bin/python
# -*- coding: utf-8 -*-
"""
Created on Fri Aug 07 17:30:58 2015
 
@author: Dreace
@modified: sndnvaps, 2016-08-02

usage:
  python pic_dl.py
"""
import urllib2
import sys
import time
import os
import random
from multiprocessing.dummy import Pool as ThreadPool 
type_ = sys.getfilesystemencoding()
def rename():
    return time.strftime("%Y%m%d%H%M%S")
def rename_2(name):  
    if len(name) == 2:  
        name = '0' + name + '.jpg' 
    elif len(name) == 1:  
        name = '00' + name + '.jpg' 
    else:  
        name = name + '.jpg' 
    return name
def download_pic(i):
    global count
    global time_out
    if Filter(i):
        try: 
            content = urllib2.urlopen(i,timeout = time_out)
            url_content = content.read()
            file_name = repr(random.randint(10000,999999999)) + "_" + rename_2(repr(count))
            f = open(file_name,"wb")
            f.write(url_content)
            f.close()
            if os.path.getsize(file_name) >= 1024*11:
                count += 1
            else:
                os.remove(file_name)
        except Exception, e:
            print e
def Filter(content):
    for line in Filter_list:
        if content.find(line) == -1:
            return True
def get_pic(url_address):
    global pic_list
    global time_out
    global headers
    try:
        req = urllib2.Request(url = url_address,headers = headers)
        str_ = urllib2.urlopen(req, timeout = time_out).read()
        url_content = str_.split("\'")
        for i in url_content:
            if i.find(".jpg") != -1:
                pic_list.append(i)   
    except Exception, e:
        print e
MAX = 100
count = 0
time_out = 60
thread_num = 50
pic_list = []
page_list = []
pic_kind = ["xinggan","share","mm","taiwan","japan"]
Filter_list = ["imgsize.ph.126.net","img.ph.126.net","img2.ph.126.net"]
dir_name = "./Photos/"+rename()
os.makedirs(dir_name)
os.chdir(dir_name)
start_time = time.time()
url_address = "http://www.mzitu.com/"
headers = {"User-Agent":" Mozilla/5.0 (Windows NT 10.0; rv:39.0) Gecko/20100101 Firefox/39.0"}
for pic_i in pic_kind:     
    for i in range(1,MAX + 1):  
        page_list.append(url_address + pic_i + "/page/" + repr(i))
page_pool = ThreadPool(thread_num)
page_pool.map(get_pic,page_list)
page_pool.close()
page_pool.join()
print "获取到".decode("utf-8").encode(type_),len(pic_list),"张图片，开始下载！".decode("utf-8").encode(type_)
pool = ThreadPool(thread_num) 
pool.map(download_pic,pic_list)
pool.close() 
pool.join()
print count,"张图片保存在".decode("utf-8").encode(type_) + dir_name
print "共耗时".decode("utf-8").encode(type_),time.time() - start_time,"s"
