#include <iostream>
#include <iomanip>
#include <fstream>
#include <sstream>
#include <vector>
#include <cmath>

#include "tbb/task_scheduler_init.h"
#include "tbb/parallel_for.h"
#include "tbb/blocked_range.h"

using namespace std;
using namespace tbb;

struct Problem {
	int id;
	int total;
	int correct;
	bool answer;
	double score;
};

class Scoring {
private:
	vector<Problem>& problems;

public:
	Scoring(vector<Problem>& problems) : problems(problems) { }

	void operator()(const blocked_range<int>& r) const {
		for (int i = r.begin(); i != r.end(); i++) {
			int total = problems[i].total + 2;
			double correct = problems[i].correct + 1;
			double prev_rate = correct / total;
			if (problems[i].answer) correct++;
			total++;
			double current_rate_inv = total / correct;
			problems[i].score *= prev_rate * current_rate_inv;
		}
	}
};

void read_problems(const char* fname, vector<Problem>& problems)
{
	fstream fs(fname, ios_base::in);
#if 0
	int id;
	int total;
	int correct;
	bool answer;
	double score;
#endif
	Problem problem;
	string line;
	stringstream ss;

	while (getline(fs, line)) {
		ss.str("");  ss.clear();
		ss.str(line);
		ss >> problem.id >> problem.total >> problem.correct >> problem.answer >> problem.score;
		problems.push_back(problem);
	}
}

int main(int argc, char* argv[])
{
	task_scheduler_init init;

	vector<Problem> problems;

	read_problems(argv[1], problems);
	parallel_for(blocked_range<int>(0, problems.size(), 8), Scoring(problems));

	for (vector<Problem>::iterator p = problems.begin(); p != problems.end(); ++p) {
		cout << p->id << " " << p->score << endl;
	}

	return 0;
}
