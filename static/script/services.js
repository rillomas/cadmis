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
	}).
	// 試験問題を取得するサービス
	factory('exam', function($http) {

		var service = {};
		service.scheme = "http";
		service.domain = "api.iknow.jp";
		service.port = "80";
		// service.itemId = "34891";

		// 項目の一覧を取得するURLを生成する
		service.generateItemsUrl = function(goalId) {
			return "{0}://{1}:{2}/goals/{3}/items?callback=JSON_CALLBACK".format(this.scheme, this.domain, this.port, goalId);
		};

		// 項目ごとの選択肢を取得するURLを取得する
		service.generateDistractorUrl = function(goalId, itemId) {
			return "{0}://{1}:{2}/goals/{3}/items/{4}/distractors?callback=JSON_CALLBACK".format(this.scheme, this.domain, this.port, goalId, itemId);
		};

		// 試験問題を取得する
		service.getExam = function (goalId, onSuccess, onError) {
			var cueId = "cue";
			var responseId = "response";
			var contentId = "content";
			var textId = "text";
			var url = service.generateItemsUrl(goalId);
			console.log("Getting exam from: " + url);
			// $http.jsonp(url).success(onSuccess).error(onError);
			$http.jsonp(url).success(function(data, status, headers, config) {

				// 問題のリストを作る
				var problemList = [];
				angular.forEach(data.items, function(item) {
					var hasCue = cueId in item;
					var hasResponse = responseId in item;

					var target = item[cueId][contentId][textId];
					var answer = item[responseId][contentId][textId];

					// 問題ごとの選択肢をつくる
					var distractorUrl = service.generateDistractorUrl(goalId, item.id);
					$http.jsonp(distractorUrl).success(function(d, s, h ,c) {
						console.log("Got distractor for " + distractorUrl);
						var optionList = [];
						var distractorList = d.distractors[item.id][responseId];
						angular.forEach(distractorList, function(opt) {
							optionList.push({
								description: opt.text,
								selected: false
							});
						});
						var problem = {
							id: item.id,
							target : target,
							answer : answer,
							optionList: optionList,
							solving: false
						};
						console.log(problem);
						problemList.push(problem);
					}).error(function(d, s, h, c) {
					});
				});
				onSuccess(problemList);	
			}).error(function(data, status, headers, config) {
				onError([]);
			});
		};

		return service;
	});
