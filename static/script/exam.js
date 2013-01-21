/**
 * 試験画面のコントローラー
 */
function ExamController($scope, $routeParams, exam) {

	$scope.goalId = "469230";
	$scope.busy = false;
	$scope.problemList = [];

	// 試験を開始する
	$scope.startExam = function(goalId) {
		$scope.busy = true;
		exam.getExam(goalId, function(problemList) {
			console.log("Get exam success");
			$scope.problemList = problemList;
			$scope.busy = false;
		}, function(problemList) {
			console.log("Get exam failed");
			console.log(data);
			$scope.busy = false;
		});
	};
}
