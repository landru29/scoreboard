angular.module("scoreboard").factory("ChronometerFactory", function ChronometerFactory ($injector, $timeout, $interval) {
    "use strict";

    /**
     * ChronometerFactory constructor
     */
    var ChronometerFactory = function ChronometerFactory (opts) {
        var service = $injector.get("Chronometer");
        this.init = 1000 * (opts.init || 0);
        this.time = moment.duration(this.init);
        service.registerChronometer(opts.id, this);
        this.countdown = !!opts.countdown;
    };

    ChronometerFactory.prototype.stop = function () {
        this.pause();
        this.time = moment.duration(this.init);
    };

    ChronometerFactory.prototype.pause = function () {
        if (this.handler) {
            $interval.cancel(this.handler);
        }
        delete this.handler
    };

    ChronometerFactory.prototype.start = function () {
        if (this.handler) {
            return;
        }
        this.handler = $interval(() => {
            if (this.countdown) {
                this.time.subtract(100);
            } else {
                this.time.add(100);
            }
        }, 100);
    };

    return ChronometerFactory;
});