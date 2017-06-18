angular.module("scoreboard").controller("GameAddCtrl", function GameAddCtrl ($scope, $state, $q, Game, Team, toaster) {
    "use strict";

    var self = this;

    /**
     * Load all teams
     * @return {Promise}
     */
    this.loadTeams = function () {
        return Team.list().$promise.then(function (teams) {
            self.teams = teams;
            self.teamPoolA = teams;
            self.teamPoolB = teams;
            return teams;
        });
    };

    this.teamSelect = function (team) {
        this.teamPoolA = _.filter(this.teams, function (team) {
            return team.id !== _.get(self.game, "teamB.id");
        });
        this.teamPoolB = _.filter(this.teams, function (team) {
            return team.id !== _.get(self.game, "teamA.id");
        });
    };

    /**
     * Add a game
     * @return {Promise}
     */
    this.addGame = function () {
        return Game.create({
            name: this.game.name,
            teamA: this.game.teamA.id,
            teamB: this.game.teamB.id
        }).$promise.then(function (game) {
            $scope.$emit("refresh-game-list");
            toaster.pop({ type: "success", title: "New game"});
            return $state.go("main.tabs.games.detail", { gameId: game.id });
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Game", body:"could not be created"});
            return $q.reject(err);
        });
    };

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.loadTeams();
    };

});