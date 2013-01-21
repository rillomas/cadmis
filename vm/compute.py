#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys, os, cgi, sqlite3, datetime, subprocess, json
import socket, fcntl
from wsgiref.simple_server import make_server

DATADIR  = "/home/futatsugi/develop/work/challengers/2013_winter/hackathon/scoring"
DATAFILE = "score.dat"

COMMAND_TRAIN = "/home/futatsugi/bin/svm-train -c 16 -g 4 -m 400 %s %s"
#COMMAND_SCORE = "setenv LD_LIBRARY_PATH /usr/local/tbb/build/linux_intel64_gcc_cc4.4.6_libc2.12_kernel2.6.32_release;/home/futatsugi/bin/scoring_tbb %s"
COMMAND_SCORE = "export LD_LIBRARY_PATH=/usr/local/tbb/build/linux_intel64_gcc_cc4.4.6_libc2.12_kernel2.6.32_release;/home/futatsugi/bin/scoring_tbb %s"
#COMMAND_SCORE = "/home/futatsugi/bin/scoring_tbb %s"

def compute(environ, start_response):
	result = "NO DATA"
	if environ["PATH_INFO"] == "/":
		if environ["REQUEST_METHOD"] == "POST":
			fs = cgi.FieldStorage(fp=environ["wsgi.input"], environ=environ, keep_blank_values=1)
			user = fs.getfirst("user", "").strip()
			data = fs.getfirst("data", "").strip()
			data = json.loads(data)
			f = file(os.path.join(DATADIR, DATAFILE), "w")
			for elems in data:
				datalist = (elems["id"], elems["total"], elems["correct"], elems["answer"], elems["score"])
				print >>f, " ".join(map(str, datalist))
			f.close()
			#os.system(COMMAND_SCORE_COPY % os.path.join(DATADIR, DATAFILE))
			#print user, data #####
			result = subprocess.Popen(COMMAND_SCORE % (os.path.join(DATADIR, DATAFILE),), shell=True, stdin=None, stdout=subprocess.PIPE, stderr=subprocess.STDOUT).stdout.read()
			result = result.strip().split("\n")
			res = []
			for r in result:
				x = r.split()
				a = {}
				a["id"] = int(x[0])
				a["score"] = float(x[1])
				res.append(a)
			result = json.dumps(res)
	#start_response("200 OK", [("Content-type", "text/html;charset=utf-8")])
	start_response("200 OK", [("Content-type", "application/json;charset=utf-8")])
	return [result.encode("utf-8")]

def main(args):
	httpd = make_server("", 8080, compute)
	httpd.serve_forever()

if __name__ == "__main__": main(sys.argv)
