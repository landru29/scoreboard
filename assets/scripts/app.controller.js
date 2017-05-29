angular.module("scoreboard").controller("MainCtrl", function MainCtrl (APP, $state) {
    "use strict";

    this.go = function (state) {
        return $state.go(state);
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.tabs = APP.tabs;

        if ($state.current.name === "main") {
            $state.go("main.teams");
        }
    };

});