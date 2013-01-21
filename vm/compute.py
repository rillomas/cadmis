#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys, os, cgi, sqlite3, datetime, subprocess, json
import socket, fcntl
from wsgiref.simple_server import make_server

COMMAND = "ssh 180.37.181.90 ls -l %s"
COMMAND_COPY = "scp %s 180.37.181.90:/home/futatsugi/develop/work/challengers/2013_winter/hackathon/svm/"
COMMAND_TRAIN = "ssh 180.37.181.90 /home/futatsugi/bin/svm-train -c 16 -g 4 -m 400 /home/futatsugi/develop/work/challengers/2013_winter/hackathon/svm/%s /home/futatsugi/develop/work/challengers/2013_winter/hackathon/svm/%s"
#COMMAND_SCORE_COPY = "scp %s 180.37.181.90:/home/futatsugi/develop/work/challengers/2013_winter/hackathon/scoring/"
#COMMAND_SCORE_COPY = "cp %s /home/futatsugi/develop/work/challengers/2013_winter/hackathon/scoring/"
#COMMAND_SCORE = 'ssh 180.37.181.90 "setenv LD_LIBRARY_PATH /usr/local/tbb/build/linux_intel64_gcc_cc4.4.6_libc2.12_kernel2.6.32_release;/home/futatsugi/bin/scoring_tbb /home/futatsugi/develop/work/challengers/2013_winter/hackathon/scoring/%s"'
#COMMAND_SCORE = "setenv LD_LIBRARY_PATH /usr/local/tbb/build/linux_intel64_gcc_cc4.4.6_libc2.12_kernel2.6.32_release;/home/futatsugi/bin/scoring_tbb %s"
COMMAND_SCORE = "export LD_LIBRARY_PATH=/usr/local/tbb/build/linux_intel64_gcc_cc4.4.6_libc2.12_kernel2.6.32_release;/home/futatsugi/bin/scoring_tbb %s"
#COMMAND_SCORE = "/home/futatsugi/bin/scoring_tbb %s"

DATAFILE = "/home/futatsugi/develop/work/challengers/2013_winter/hackathon/scoring/score.dat"

def compute(environ, start_response):
	result = ""
	if environ["PATH_INFO"] == "/":
		#result = subprocess.Popen(COMMAND % (".",), shell=True, stdin=None, stdout=subprocess.PIPE, stderr=subprocess.STDOUT).stdout.read()
		#result = result.replace("\n", "<br />")
		if environ["REQUEST_METHOD"] == "POST":
			fs = cgi.FieldStorage(fp=environ["wsgi.input"], environ=environ, keep_blank_values=1)
			user = fs.getfirst("user", "").strip()
			data = fs.getfirst("data", "").strip()
			data = json.loads(data)
			f = file(DATAFILE, "w")
			#for id, elems in data.iteritems():
			for elems in data:
				datalist = (elems["id"], elems["total"], elems["correct"], elems["answer"], elems["score"])
				#datalist = (id, elems["total"], elems["correct"], elems["answer"], elems["score"])
				"""
				id = elems["id"]
				total = elems["total"]
				correct = elems["correct"]
				answer = elems["answer"]
				score = elems["score"]
				"""
				#f.write(" ".join(map(str, datalist)))
				print >>f, " ".join(map(str, datalist))
			f.close()
			#os.system(COMMAND_SCORE_COPY % DATAFILE)
			#print user, data #####
			result = subprocess.Popen(COMMAND_SCORE % (DATAFILE,), shell=True, stdin=None, stdout=subprocess.PIPE, stderr=subprocess.STDOUT).stdout.read()
			result = result.strip().split("\n")
			res = []
			for r in result:
				#print r #####
				x = r.split()
				a = {}
				a["id"] = int(x[0])
				a["score"] = float(x[1])
				res.append(a)
			#result = result.replace("\n", "<br />")
			result = json.dumps(res)
	#start_response("200 OK", [("Content-type", "text/html;charset=utf-8")])
	start_response("200 OK", [("Content-type", "application/json;charset=utf-8")])
	return [result.encode("utf-8")]

def main(args):
	httpd = make_server("", 8080, compute)
	httpd.serve_forever()

if __name__ == "__main__": main(sys.argv)
