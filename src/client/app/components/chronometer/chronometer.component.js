angular.module("scoreboard").component("chronometer", {
    bindings: {
        countdown: "=?",
        init: "@",
        format: "@",
        id: "@"
    },
    templateUrl: "app/components/chronometer/chronometer.html",
    controller: function chronometerCtrl (ChronometerFactory) {
        "use strict";

        /**
        * Initialization of the component
        */
        this.$onInit = function () {
            this.format = this.format || "HH:mm";
            this.chrono = new ChronometerFactory({
                init: this.init,
                id: this.id
            });
        };
    }
});