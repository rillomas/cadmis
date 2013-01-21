/**
 * 試験画面のコントローラー
 */
function ExamController($scope, $routeParams, exam, authenticate) {

	$scope.busy = false;
	$scope.exam = null;
	$scope.startTime = null;

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
		exam.submitExam(examResult, $scope.startTime, endTime, parseInt(authenticate.userId,10), function() {
			$scope.busy = false;
		}, function() {
			$scope.busy = false;
		});
	};
}
