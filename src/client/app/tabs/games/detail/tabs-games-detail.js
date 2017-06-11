angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.tabs.games.detail",
        url: "/:gameId",
        views: {
            gameDetail: {
                templateUrl: "app/tabs/games/detail/tabs-games-detail.html",
                controller: "GameDetailCtrl",
                controllerAs: "GameDetailCtrl"
            }
        }
    });
});