angular.module("scoreboard").filter("chronometer", function digits () {
    "use strict";

    function resize (input, size) {
        var num = parseInt("" + (input || "0"), 10);
        return (new Array(size + 1).join("0") + num).substr(-size);
    }

    return function digits (input) {
        if (!input || typeof input !== "object") {
            return "00:00"
        }
        var result = [resize(input.minutes(), 2), resize(input.seconds(), 2)];
        if (input.hours()) {
            result.unshift(input.hours());
        }
        return result.join(":");
    };
});