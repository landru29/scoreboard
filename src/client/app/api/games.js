angular.module("scoreboard").service("Game", function Team ($resource) {
    "use strict";

    return $resource(
        "/games",
        {
            gameId: "@gameId",
        }, {
            list: {
                url: "/games",
                method: "GET",
                isArray: true
            },
            detail: {
                url: "/games/:gameId",
                method: "GET",
                isArray: false
            },
            delete: {
                url: "/games/:gameId",
                method: "DELETE",
                isArray: false
            },
            update: {
                url: "/games/:gameId",
                method: "PUT",
                isArray: false
            },
            create: {
                url: "/games",
                method: "POST",
                isArray: false
            }
        }
        );
});