/**
 * プロフィール画面のコントローラー
 */
function ProfileController($scope, $routeParams, user) {
	var accessToken = sessionStorage.getItem(constants.AccessTokenKey);
	var userId = sessionStorage.getItem(constants.UserIdKey);
	$scope.busy = false;
	$scope.errorMessage = '';
	$scope.profile = user.getProfile(userId, accessToken);
	$scope.apply = function(profile){
		$scope.busy = true;
		user.applyProfile(profile, function () {
			console.log("Apply profile success");
			$scope.busy = false;
		}, function() {
			console.log("Apply rpfile error");
			$scope.busy = false;
		});
	};
}
