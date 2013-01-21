/**
 * Services
 */
angular.module('cadmis.service', ['ngResource']).
	// ユーザー関連のリソースを管理するサービス
	factory('user', function($resource) {
		var service = {};

		// 新規ユーザー追加機能
		service.signUp = function(email, password, onSuccess, onError) {
			var User = $resource('/api/1/user');
			var newUser = new User();
			newUser.Email = email;
			newUser.Password = password;
			newUser.$save({}, onSuccess, onError);
		};

		// プロフィールを取得する
		service.getProfile = function(userId, accessToken) {
			var Profile = $resource('/api/1/user_profile', { ui: userId, at: accessToken});
			var profile = Profile.get( {}, function() {
			});
			return profile;
		};

		// プロフィールを適用する
		service.applyProfile = function(profile, onSuccess, onError) {
			var Profile = $resource('/api/1/user_profile');
			var newProf = new Profile();
			newProf.Profile = profile;
			newProf.$save({}, onSuccess, onError);
		};
		return service;
	}).
	// 認証トークンを管理するサービス
	factory('authenticate', function($resource, $rootScope) {
		var service = {};

		// session storageに保存したアクセストークンがあればそれを使う
		var storage = sessionStorage;
		var token = storage.getItem(constants.AccessTokenKey);
		var userId = storage.getItem(constants.UserIdKey);

		// ログイン済みかを示すトークン
		service.accessToken = token;

		// ユーザーID
		service.userId = userId; 

		// ログイン済みかどうか
		service.authenticated = function() {
			return this.accessToken != null;
		};

		service.notifyAuthenticationChanged = function() {
			var args = { "authenticated" : this.authenticated() };
			$rootScope.$emit(constants.AuthenticationChangedEvent, args);
		};

		// トークンをリクエストする
		service.requestToken = function(email, password, onSuccess, onError) {
			var AccessToken = $resource('/api/1/access_token');
			var newToken = new AccessToken();
			newToken.Email = email;
			newToken.Password = password;

			var success = function(response) {
				console.log("Got access token: " + response.Id);
				var token = response.Id;
				var userId = response.UserId;
				service.accessToken = token;
				service.userId = userId;

				// ブラウザ側にトークンを保存する
				sessionStorage.setItem(constants.AccessTokenKey, service.accessToken);
				sessionStorage.setItem(constants.UserIdKey, service.userId);

				onSuccess(response);

				// 認証完了したのでイベントを飛ばす
				service.notifyAuthenticationChanged();
			};
			newToken.$save({}, success, onError);
		};

		// トークンを破棄する
		service.disposeToken = function() {
			this.accessToken = null;
			sessionStorage.clear();
			this.notifyAuthenticationChanged();
		};

		return service;
	});
