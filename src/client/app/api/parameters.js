angular.module("scoreboard").service("Parameters", function Team ($resource) {
    "use strict";

    return $resource(
        "/parameters",
        {}, {
            read: {
                url: "/parameters",
                method: "GET",
                isArray: false
            },
            update: {
                url: "/parameters",
                method: "POST",
                isArray: false
            }
        }
        );
});