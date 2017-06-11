angular.module("scoreboard").controller("GameDetailCtrl", function GameDetailCtrl ($stateParams, $state, $scope, $q, Team, Game, toaster) {
    "use strict";
    var self = this;

    /**
     * Load all the teams
     * @return {Promise}
     */
    this.loadTeams = function () {
        return Team.list().$promise.then(function (teams) {
            self.teams = teams;
            return teams;
        });
    };

    /**
     * Get the detail of a team
     * @return {Promise}
     */
    this.getGameDetail = function () {
        return Game.detail({
            gameId: $stateParams.gameId
        }).$promise.then(function (game) {
            self.game = game;
            self.gameTitle = game.name;
            self.game.teamA = _.find(self.teams, { id: self.game.teamA.id });
            self.game.teamB = _.find(self.teams, { id: self.game.teamB.id });
            return game;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Game", body:"could not be found"});
            return $state.go("main.tabs.games");
        });
    };

    /**
     * Delete the current game
     * @return {Promise}
     */
    this.deleteGame = function () {
        this.gameBusy = true;
        return Game.delete({
            gameId: $stateParams.gameId
        }).$promise.then(function () {
            $scope.$emit("refresh-game-list");
            toaster.pop({ type: "success", title: "Game deleted"});
            return $state.go("main.tabs.games");
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Game", body:"could not be deleted"});
            return $q.reject(err);
        }).finally(function() {
            delete self.gameBusy;
        });
    };

    /**
     * Update the current game
     * @return {Promise}
     */
    this.saveGame = function () {
        this.gameBusy = true;
        return Game.update({
            gameId: $stateParams.gameId
        }, {
            name: this.game.name,
            teamA: this.game.teamA.id,
            teamB: this.game.teamB.id
        }).$promise.then(function (game) {
            $scope.$emit("refresh-game-list");
            self.gameTitle = game.name;
            toaster.pop({ type: "success", title: "Game saved"});
            return game;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Game", body:"could not be saved"});
            return $q.reject(err);
        }).finally(function() {
            delete self.gameBusy;
        });
    };

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.loadTeams().then(function () {
            return self.getGameDetail();
        });
    };

});