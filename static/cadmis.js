var constants = {
	"AuthenticationChangedEvent" : "AuthenticationChangedEvent"
};

/**
 * Services
 */
angular.module('cadmis.service', ['ngResource']).
	// ユーザー関連のリソースを管理するサービス
	factory('user', function($resource) {
		var service = {};

		// 新規ユーザー追加機能
		service.signUp = function(email, password, onSuccess, onError) {
			var User = $resource('/api/1/user/');
			var newUser = new User();
			newUser.Email = email;
			newUser.Password = password;
			newUser.$save({}, onSuccess, onError);
		}
		return service;
	}).
	// 認証トークンを管理するサービス
	factory('authenticate', function($resource, $rootScope) {
		var service = {};

		// ログイン済みかを示すトークン
		service.accessToken = null;

		// トークンをリクエストする
		service.requestToken = function(email, password, onSuccess, onError) {
			var AccessToken = $resource('/api/1/access_token');
			var newToken = new AccessToken();
			newToken.Email = email;
			newToken.Password = password;

			var success = function(response) {
				console.log("Got access token");
				this.accessToken = response.data;
				onSuccess(response);
				
				// 認証完了したのでイベントを飛ばす
				var args = { "authenticated" : true };
				$rootScope.$emit(constants.AuthenticationChangedEvent, args);
			}
			newToken.$save({}, success, onError);
		}

		// ログイン済みかどうか
		service.authenticated = function() {
			return this.aceessToken != null;
		}
		return service;
	});

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

				$scope.signUp = function () {
					var email = $scope.email;
					var pass = $scope.password;
					console.log(
						"Email: " + email +
						" Password: " + pass);

					user.signUp(email, pass, function() {
						console.log("sign up success");
						$scope.errorMessage = '';

						// ログインする
						// authenticate.requestToken(email, pass, function() {
						// 	console.log("login success");
						// }, function(response) {
						// 	console.log("login error");
						// 	console.log(response.data);
						// 	$scope.errorMessage = response.data;
						// });
					}, function(response) {
						console.log("sign up error");
						console.log(response.data);
						$scope.errorMessage = response.data;
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

				// send id/pass to server and authenticate
				$scope.login = function (email, password, rememberLogin) {
					console.log(
						"email: " + email +
						" Password: " + password +
						" Remember: " + rememberLogin);

					authenticate.requestToken(email, password, function(response) {
					}, function(response) {
					})
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
function CadmisController($scope, $rootScope) {
    $scope.authenticated = false;

    $rootScope.$on(constants.AuthenticationChangedEvent, function(event, args) {
    	console.log("Got AuthenticationChangedEvent: " + args.authenticated);
    	$scope.authenticated = args.authenticated;
    });
}