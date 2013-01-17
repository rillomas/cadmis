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

		// ログイン済みかどうか
		service.authenticated = function() {
			return this.accessToken != null;
		}

		service.notifyAuthenticationChanged = function() {
			var args = { "authenticated" : this.authenticated() };
			$rootScope.$emit(constants.AuthenticationChangedEvent, args);
		}

		// トークンをリクエストする
		service.requestToken = function(email, password, onSuccess, onError) {
			var AccessToken = $resource('/api/1/access_token');
			var newToken = new AccessToken();
			newToken.Email = email;
			newToken.Password = password;

			var success = function(response) {
				console.log("Got access token: " + response);
				service.accessToken = response[0];
				onSuccess(response);

				// 認証完了したのでイベントを飛ばす
				service.notifyAuthenticationChanged();
			}
			newToken.$save({}, success, onError);
		}

		// トークンを破棄する
		service.disposeToken = function() {
			this.accessToken = null;
			this.notifyAuthenticationChanged();
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
function CadmisController($scope, $rootScope, authenticate) {
    $scope.authenticated = false;

    $scope.logout = function () {
    	authenticate.disposeToken();
    }

    $rootScope.$on(constants.AuthenticationChangedEvent, function(event, args) {
    	console.log("Got AuthenticationChangedEvent: " + args.authenticated);
    	$scope.authenticated = args.authenticated;
    });
}