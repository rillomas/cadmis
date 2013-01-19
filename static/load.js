"use strict";

/**
 * 必要なスクリプトを読み込む
 */
head.js(
	"http://ajax.googleapis.com/ajax/libs/jquery/1.8.3/jquery.min.js",
	"http://ajax.googleapis.com/ajax/libs/angularjs/1.0.3/angular.min.js",
	"http://ajax.googleapis.com/ajax/libs/angularjs/1.0.3/angular-resource.min.js",
	"script/bootstrap.min.js",

	// アプリ全体で使用する部品
	"script/constants.js",
	"script/services.js",
	"script/components.js",

	// 各画面ごとのコントローラー
	"script/userHome.js",
	"script/ranking.js",
	"script/exam.js",
	"script/profile.js",

	// アプリそのもののコントローラー
	"script/cadmis.js"
	);
