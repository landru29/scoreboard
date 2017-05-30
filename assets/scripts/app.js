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


(function (context) {
    context.onload = function () {
        var conn;

        if (context.WebSocket) {
            conn = new WebSocket("ws://" + document.location.host + "/ws/");
            conn.onclose = function (evt) {
                console.log("Websocket connection closed");
            };
            conn.onmessage = function (evt) {
                _.forEach(evt.data.split('\n'), function (line) {
                    console.log("WS", line);
                });
            };
        } else {
            console.warn("Your browser does not support websockets");
        }

        context.sendWs = function (str) {
            if (!conn) {
                return false;
            }
            if (!str) {
                return false;
            }
            conn.send(str);
            return false;
        };
    };
})(window);