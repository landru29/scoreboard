angular.module("scoreboard").controller("MainCtrl", function MainCtrl ($scope, $state) {
    "use strict";

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        if ($state.current.name === "main") {
            $state.go("main.tabs.board");
        }
    };

});