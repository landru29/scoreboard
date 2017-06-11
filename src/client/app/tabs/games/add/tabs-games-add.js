angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.tabs.games.add",
        url: "/add",
        views: {
            gameDetail: {
                templateUrl: "app/tabs/games/add/tabs-games-add.html",
                controller: "GameAddCtrl",
                controllerAs: "GameAddCtrl"
            }
        }
    });
});