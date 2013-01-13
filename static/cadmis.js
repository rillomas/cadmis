/**
 *  Components
 */
angular.module('cadmis.component',[]).
	directive('signupForm', function() {
		return {
			restrict: 'E',
			transclude: false,
			scope: {},
			controller: function ($scope, $element) {

				$scope.userId = '';
				$scope.password = '';

				$scope.signUp = function () {
					console.log(
						"Id: " + $scope.userId +
						" Password: " + $scope.password);
				};
			},
			templateUrl: 'component/signupForm.html',
			replace: true
		};
	}).
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