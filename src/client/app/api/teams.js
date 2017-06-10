angular.module("scoreboard").service("Team", function Team ($resource) {
    "use strict";

    return $resource(
        "/teams/:teamId",
        {
            teamId: "@teamId"
        }, {
            list: {
                url: "/teams",
                method: "GET",
                isArray: true
            },
            detail: {
                url: "/teams/:teamId",
                method: "GET",
                isArray: false
            },
            delete: {
                url: "/teams/:teamId",
                method: "DELETE",
                isArray: false
            },
            update: {
                url: "/teams/:teamId",
                method: "PUT",
                isArray: false
            },
            create: {
                url: "/teams",
                method: "POST",
                isArray: false
            }
        }
        );
});