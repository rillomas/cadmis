# Hackathon 2013.01.21-22 by N. Futatsugi

CC = g++

FLAGS += -O3
FLAGS += -ltbb
FLAGS += -L /usr/local/tbb/build//usr/local/tbb/build/linux_intel64_gcc_cc4.4.5_libc2.11.3_kernel2.6.32_release
FLAGS += -I /usr/local/tbb/include

all:
	$(CC) $(FLAGS) scoring_tbb.cpp -o scoring_tbb
	cp scoring_tbb ~/bin/
