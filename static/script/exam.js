/**
 * 試験画面のコントローラー
 */
function ExamController($scope, $routeParams, exam, authenticate) {

	$scope.busy = false;
	$scope.exam = null;
	$scope.startTime = null;
	$scope.result = null;
	$scope.repeating = false;

	// 試験を開始する
	$scope.startExam = function() {
		$scope.busy = true;
		exam.getExam(function(exam) {
			console.log("Get exam success");
			console.log(exam);
			$scope.exam = exam;
			$scope.busy = false;
			$scope.startTime = new Date().getTime();
		}, function(exam) {
			console.log("Get exam failed");
			$scope.busy = false;
		});
	};

	// 試験問題を提出する
	$scope.submitExam = function(examResult) {
		var endTime = new Date().getTime();
		$scope.busy = true;
		exam.submitExam(examResult, $scope.startTime, endTime, parseInt(authenticate.userId,10), function(result) {
			$scope.busy = false;
			$scope.result = result;

			var correctNum = 0;
			var totalNum = 0;
			angular.forEach(result.ProblemList, function(problem) {
				totalNum++;
				if (problem.Correct) {
					correctNum++;
				}
			});
			$scope.resultMessage = "Your test result was {0}/{1}".format(correctNum, totalNum);
		}, function() {
			$scope.busy = false;
		});
	};

	$scope.startNextExam = function() {
		$scope.repeating = true;
		$scope.result = null;
		$scope.resultMessage = null;
		$scope.exam = null;
		$scope.startExam();
	};
}
