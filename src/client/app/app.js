angular.module("scoreboard", [
    "ui.router",
    "scoreboard.templates",
    "ui.bootstrap",
    "ngResource",
    "ui.bootstrap.materialPicker",
    "angularFileUpload",
    "toaster",
    "ngAnimate",
    "ui.select",
    "ngSanitize"
]);

angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main",
        url: "",
        views: {
            mainView: {
                templateUrl: "app/app.html",
                controller: "MainCtrl",
                controllerAs: "MainCtrl"
            }
        }
    });
});

angular.module("scoreboard").config(function(WsProvider) {
    WsProvider.autoReconnect = true;
    WsProvider.createConnection("operator", "/ws/")
});

angular.module("scoreboard").run(function($rootScope, Ws) {
    Ws.registerScope($rootScope);
});