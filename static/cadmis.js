/**
 *  Components
 */
angular.module('cadmis.component',['ngResource']).
    // ユーザー登録用のフォーム
	directive('signupForm', function() {
		return {
			restrict: 'E',
			transclude: false,
			scope: {},
			controller: function ($scope, $element, $resource) {

				$scope.email = '';
				$scope.password = '';
				$scope.errorMessage = '';

				$scope.signUp = function () {

					var email = $scope.email;
					var pass = $scope.password;
					console.log(
						"Email: " + email +
						" Password: " + pass);

					var User = $resource('/api/1/user/');

					var newUser = new User();
					newUser.Email = email;
					newUser.Password = pass;
					newUser.$save({}, function() {
						console.log("sign up success");
					}, function(data, headers) {
						console.log("sign up error");
						$scope.errorMessage = "Specified Email address is already used.";
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
			controller: function ($scope, $element) {

				$scope.loginId = '';
				$scope.password = '';
				$scope.rememberLogin = false;

				// send id/pass to server and authenticate
				$scope.authenticate = function () {
					console.log(
						"Id: " + $scope.loginId +
						" Password: " + $scope.password +
						" Remember: " + $scope.rememberLogin);
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
	}).
	// ログインされた状態で表示される画面
	directive('userHome', function() {
		return {
			restrict: 'E',
			transclude: false,
			scope: {},
			controller: function ($scope, $element) {

			},
			templateUrl: 'component/userHome.html',
			replace: true
		}
	});

/**
 * Main module
 */
angular.module('cadmis',['cadmis.component']).
    run(function () {
    });

/**
 * Main controller
 */
function CadmisController($scope) {
    $scope.authenticated = false;
}