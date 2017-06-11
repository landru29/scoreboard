angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main.tabs",
        url: "/tabs",
        views: {
            mainContent: {
                templateUrl: "app/tabs/tabs.html",
                controller: "TabsCtrl",
                controllerAs: "TabsCtrl"
            }
        }
    });
});