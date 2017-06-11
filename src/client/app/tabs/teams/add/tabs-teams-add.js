angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.tabs.teams.add",
        url: "/add",
        views: {
            teamDetail: {
                templateUrl: "app/tabs/teams/add/tabs-teams-add.html",
                controller: "TeamAddCtrl",
                controllerAs: "TeamAddCtrl"
            }
        }
    });
});