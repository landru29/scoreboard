angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.teams.detail",
        url: "/:teamId",
        views: {
            teamDetail: {
                templateUrl: "app/tabs/teams/detail/tabs-teams-detail.html",
                controller: "TeamDetailCtrl",
                controllerAs: "TeamDetailCtrl"
            }
        }
    });
});