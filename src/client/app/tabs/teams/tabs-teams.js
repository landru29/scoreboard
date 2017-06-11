angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.tabs.teams",
        url: "/teams",
        views: {
            tabContent: {
                templateUrl: "app/tabs/teams/tabs-teams.html",
                controller: "TeamsCtrl",
                controllerAs: "TeamsCtrl"
            }
        }
    });
});