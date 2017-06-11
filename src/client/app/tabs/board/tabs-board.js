angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.tabs.board",
        url: "/board",
        views: {
            tabContent: {
                templateUrl: "app/tabs/board/tabs-board.html",
                controller: "BoardCtrl",
                controllerAs: "BoardCtrl"
            }
        }
    });
});