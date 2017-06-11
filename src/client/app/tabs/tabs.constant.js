angular.module("scoreboard").constant("TABS", {
    tabs: [
        {
            name: "Board",
            icon: "fa fa-eye",
            state: "main.tabs.board"
        }, {
            name: "Teams",
            icon: "fa fa-users",
            state: "main.tabs.teams"
        }, {
            name: "Games",
            icon: "fa fa-trophy",
            state: "main.tabs.games"
        }
    ]
});