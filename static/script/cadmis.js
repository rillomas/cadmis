/**
 * Main controller
 */
function CadmisController($scope, $rootScope, $location, authenticate) {
    $scope.authenticated = authenticate.authenticated();
    $scope.location = $location; // htmlからURLを参照できるようにlocationサービスを保持する

    $scope.logout = function () {
    	authenticate.disposeToken();
    }

    // 認証状態が変わったら同期する
    $rootScope.$on(constants.AuthenticationChangedEvent, function(event, args) {
    	console.log("Got AuthenticationChangedEvent: " + args.authenticated);
    	$scope.authenticated = args.authenticated;
    });
}

/**
 * Main module
 */
angular.module('cadmis',['cadmis.component']).
	config(function($routeProvider, $locationProvider) {
		// アプリ全体の設定
		
		// html5モードを使う
		$locationProvider.html5Mode(true);

		// URLごとにViewとコントローラーを割り当てる
		$routeProvider.when('/ranking', {
			templateUrl : 'component/ranking.html',
			controller : RankingController,
		});
		$routeProvider.when('/exam', {
			templateUrl : 'component/exam.html',
			controller : ExamController,
		});	
		$routeProvider.when('/profile', {
			templateUrl : 'component/profile.html',
			controller : ProfileController,
		});
		$routeProvider.otherwise( {
			templateUrl : 'component/userHome.html',
			controller : UserHomeController,
		})
	});
