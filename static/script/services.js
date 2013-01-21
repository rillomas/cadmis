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
	factory('exam', function($resource, $http) {

		var service = {};
		service.scheme = "http";
		service.domain = "api.iknow.jp";
		service.port = "80";
		service.examId = "469230";

		// service.itemId = "34891";

		// 項目の一覧を取得するURLを生成する
		service.generateItemsUrl = function(examId) {
			return "{0}://{1}:{2}/goals/{3}/items?callback=JSON_CALLBACK".format(this.scheme, this.domain, this.port, examId);
		};

		// 項目ごとの選択肢を取得するURLを取得する
		service.generateDistractorUrl = function(examId, itemId) {
			return "{0}://{1}:{2}/goals/{3}/items/{4}/distractors?callback=JSON_CALLBACK".format(this.scheme, this.domain, this.port, examId, itemId);
		};

		// 試験問題を取得する
		service.getExam = function (onSuccess, onError) {
			var examId = service.examId;
			var cueId = "cue";
			var responseId = "response";
			var contentId = "content";
			var textId = "text";
			var url = service.generateItemsUrl(examId);
			console.log("Getting exam from: " + url);
			$http.jsonp(url).success(function(data, status, headers, config) {

				var examOutput = { id: examId, problemList: [] };
				// 問題のリストを作る
				var maxProblemNum = 10; // とりあえず最大10個
				var dataList = data.items.slice(0, maxProblemNum);
				var problemList = [];
				angular.forEach(dataList, function(item) {
					var target = item[cueId][contentId][textId];
					var answer = item[responseId][contentId][textId];

					// console.log("target: "+ target + " answer: " + answer);

					// 問題ごとの選択肢をつくる
					var distractorUrl = service.generateDistractorUrl(examId, item.id);
					$http.jsonp(distractorUrl).success(function(d, s, h ,c) {
						// console.log("Got distractor for " + distractorUrl);
						var optionList = [];
						var distractorList = d.distractors[item.id][responseId];
						angular.forEach(distractorList, function(opt) {
							optionList.push({
								description: opt.text
							});
						});

						// 答えを選択肢に含める
						var index = Math.floor(Math.random() * optionList.length);
						// console.log("Inserting answer to index: " + index);
						var ans = { description: answer};
						optionList.splice(index, 0, ans);
						var problem = {
							id: item.id,
							target : target,
							answer : answer,
							optionList: optionList,
							solving: false,
							selected : optionList[0].description
						};
						// console.log(problem);
						problemList.push(problem);
					}).error(function(d, s, h, c) {
					});
				});
				examOutput.problemList = problemList;
				onSuccess(examOutput);	
			}).error(function(data, status, headers, config) {
				onError({});
			});
		};

		// 問題を提出する
		service.submitExam = function(examResult, startTime, finishedTime, userId, onSuccess, onError) {
			var result = {
				UserId: userId,
				ExamId: parseInt(examResult.id, 10),
				StartTime: startTime,
				FinishedTime: finishedTime,
				ProblemList: []
			};

			angular.forEach(examResult.problemList, function(problem) {
				var correct = problem.selected == problem.answer;
				var pr = {
					ProblemId: problem.id,
					Correct: correct
				};
				result.ProblemList.push(pr);
			});

			var Exam = $resource('/api/1/exam');
			var newExam = new Exam();
			newExam.Result = result;
			newExam.$save({}, function() {
				onSuccess(result);
			}, function() {
				onError(result);
			});
		};

		return service;
	}).
  // ランキング取得サービス
  factory('ranking', function($resource, $http) {
    
    var service = {};
    
    service.getUsers = function(id, handler) {
      $http.get('/api/1/compute_rank?examId={0}'.format(id)).success(handler);
    }

    service.getGoals = function(id, handler) {
      $http.get('/api/1/compute_rank?userId={0}'.format(id)).success(handler);
    }
    
    return service;
  }
  );
