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

		// session storageに保存したアクセストークンがあればそれを使う
		var storage = sessionStorage;
		var token = storage.getItem(constants.AccessTokenKey);

		// ログイン済みかを示すトークン
		service.accessToken = token;

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
				console.log("Got access token: " + response.Id);
				var token = response.Id;
				service.accessToken = token;

				// ブラウザ側にトークンを保存する
				sessionStorage.setItem(constants.AccessTokenKey, service.accessToken);

				onSuccess(response);

				// 認証完了したのでイベントを飛ばす
				service.notifyAuthenticationChanged();
			}
			newToken.$save({}, success, onError);
		}

		// トークンを破棄する
		service.disposeToken = function() {
			this.accessToken = null;
			sessionStorage.removeItem(constants.AccessTokenKey);
			this.notifyAuthenticationChanged();
		}

		return service;
	});