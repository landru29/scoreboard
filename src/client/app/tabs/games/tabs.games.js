angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.tabs.games",
        url: "/game",
        views: {
            tabContent: {
                templateUrl: "app/tabs/games/tabs-games.html",
                controller: "GameCtrl",
                controllerAs: "GameCtrl"
            }
        }
    });
});