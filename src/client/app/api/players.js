angular.module("scoreboard").service("Player", function Team ($resource) {
    "use strict";

    return $resource(
        "/teams/:teamId/players",
        {
            teamId: "@teamId",
            playerId: "@playerId"
        }, {
            list: {
                url: "/teams/:teamId/players",
                method: "GET",
                isArray: true
            },
            detail: {
                url: "/teams/:teamId/players/:playerId",
                method: "GET",
                isArray: false
            },
            delete: {
                url: "/teams/:teamId/players/:playerId",
                method: "DELETE",
                isArray: false
            },
            update: {
                url: "/teams/:teamId/players/:playerId",
                method: "PUT",
                isArray: false
            },
            create: {
                url: "/teams/:teamId/players",
                method: "POST",
                isArray: false
            }
        }
        );
});