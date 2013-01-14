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

				$scope.userId = '';
				$scope.password = '';

				$scope.signUp = function () {

					var id = $scope.userId;
					var pass = $scope.password;
					console.log(
						"UserId: " + id +
						" Password: " + pass);

					var User = $resource('/api/1/user/:UserId', {UserId: '@id'});

					var newUser = new User();
					newUser.UserId = id;
					newUser.Password = pass;
					newUser.$save();

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