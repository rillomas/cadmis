/**
 *  Components
 */
angular.module('cadmis.component',['cadmis.service']).
    // ユーザー登録用のフォーム
	directive('signupForm', function() {
		return {
			restrict: 'E',
			transclude: false,
			scope: {},
			controller: function ($scope, $element, user, authenticate) {

				$scope.email = '';
				$scope.password = '';
				$scope.errorMessage = '';
				$scope.busy = false;

				$scope.signUp = function () {
					$scope.busy = true;
					var email = $scope.email;
					var pass = $scope.password;
					console.log(
						"Email: " + email +
						" Password: " + pass);

					user.signUp(email, pass, function() {
						$scope.errorMessage = '';

						// ログインする
						authenticate.requestToken(email, pass, function(response) {
							$scope.errorMessage = '';
							$scope.busy = false;
						}, function(response) {
							$scope.errorMessage = response.data;
							$scope.busy = false;
						});
					}, function(response) {
						$scope.errorMessage = response.data;
						$scope.busy = false;
					});
				};
			},
			templateUrl: 'component/signupForm.html',
			replace: true
		};
	}).
    // ログイン用のフォーム
	directive('loginForm', function() {
		return {
			restrict: 'E',
			transclude: false,
			scope: {},
			controller: function ($scope, $element, authenticate) {

				$scope.email = '';
				$scope.password = '';
				$scope.rememberLogin = false;
				$scope.errorMessage = '';
				$scope.busy = false;

				// send id/pass to server and authenticate
				$scope.login = function (email, password, rememberLogin) {
					$scope.busy = true;
					console.log(
						"email: " + email +
						" Password: " + password +
						" Remember: " + rememberLogin);

					authenticate.requestToken(email, password, function(response) {
						$scope.errorMessage = '';
						$scope.busy = false;
					}, function(response) {
						$scope.errorMessage = response.data;
						$scope.busy = false;
					});
				};
			},
			templateUrl: 'component/loginForm.html',
			replace: true
		};
	}).
	// ログインしてない状態で表示される画面
	directive('greeting', function() {
		return {
			restrict: 'E',
			transclude: false,
			scope: {},
			controller: function ($scope, $element) {
			},
			templateUrl: 'component/greeting.html',
			replace: true
		};
	});
