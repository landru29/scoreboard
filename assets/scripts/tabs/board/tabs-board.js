angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.board",
        url: "/board",
        views: {
            tabContent: {
                templateUrl: "/scoreboard/scripts/tabs/board/tabs-board.html",
                controller: "BoardCtrl",
                controllerAs: "BoardCtrl"
            }
        }
    });
});