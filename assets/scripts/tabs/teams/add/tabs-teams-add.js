angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.teams.add",
        url: "/add",
        views: {
            teamDetail: {
                templateUrl: "/scoreboard/scripts/tabs/teams/add/tabs-teams-add.html",
                controller: "TeamAddCtrl",
                controllerAs: "TeamAddCtrl"
            }
        }
    });
});