angular.module("scoreboard").factory("ChronometerFactory", function ChronometerFactory ($injector) {
    "use strict";

    /**
     * ChronometerFactory constructor
     */
    var ChronometerFactory = function ChronometerFactory (opts) {
        var service = $injector.get("Chronometer")
        this.time = moment.duration(1000 * (opts.init || 0));
        service.registerChronometer(opts.id, this);
    };

    return ChronometerFactory;
});