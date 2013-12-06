(function () {
	"use strict";

	angular.
		module("jobmon", [
			"ngRoute"
		]).
		config(
			function (
				$logProvider,
				$locationProvider,
				$routeProvider
			) {
				$logProvider.debugEnabled(true);
				$locationProvider.html5Mode(true);

				$routeProvider.
					when("/", {
						templateUrl: "/assets/jobs_index.html",
						controller: "ListCtrl"
					}).
					when("/logs/:id", {
						templateUrl: "/assets/logs_show.html",
						controller: "LogCtrl"
					}).
					otherwise({
						redirectTo: "/"
					})
		}).
		controller("LogCtrl",
			function (
				$scope,
				$http,
				$routeParams,
				$log
			) {
				var id = $routeParams.id;
				$http.get("/logs/" + id + ".json").
					then(function (response) {
						$scope.log = response.data;
					}, function (response) {
						$log.error(response);
					});
			}).
		controller("ListCtrl",
			function (
				$scope,
				$http,
				$log
			)  {
				$http.get("/jobs.json").
					then(function (response) {
						$scope.jobs = response.data;
						$log.debug(response.data);
					}, function (response) {
						$log.error(response);
					});
			})
})()
