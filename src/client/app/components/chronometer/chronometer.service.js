angular.module("scoreboard").service("Chronometer", function Service () {
    "use strict";

    var Service = function Service () {
        this.chronometers = {};
    };

    Service.prototype.registerChronometer = function registerChronometer (id, chronometer) {
        this.chronometers[id] = chronometer;
    };

    return new Service();
});