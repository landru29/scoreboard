angular.module("scoreboard", [
    "ui.router",
    "ui.bootstrap",
    "ngResource",
    "ui.bootstrap.materialPicker",
    "angularFileUpload",
    "toaster",
    "ngAnimate"
]);

angular.module("scoreboard").config(function($stateProvider) {
    $stateProvider.state({
        name: "main",
        url: "",
        views: {
            mainView: {
                templateUrl: "/scoreboard/scripts/app.html",
                controller: "MainCtrl",
                controllerAs: "MainCtrl"
            }
        }
    });
});