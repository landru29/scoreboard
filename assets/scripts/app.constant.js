angular.module("scoreboard").constant("APP", {
    tabs: [
        {
            name: "Teams",
            icon: "fa fa-users",
            state: "main.teams"
        }, {
            name: "Board",
            icon: "fa fa-eye",
            state: "main.board"
        }
    ]
});