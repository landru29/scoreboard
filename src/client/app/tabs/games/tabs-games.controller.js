angular.module("scoreboard").controller("GameCtrl", function GameCtrl ($q, $scope, Game, Parameters, toaster) {
    "use strict";

    var self = this;

    /**
     * Load all teams
     * @return {Promise}
     */
    this.loadGames = function () {
        return Game.list().$promise.then(function (games) {
            self.games = games;
            return games;
        });
    }

    /**
     * Load the current game parameters
     * @return {Promise}
     */
    this.loadGameParameters = function () {
        return Parameters.read().$promise.then(function (parameters) {
            self.parameters = parameters;
            self.parameters.game = _.find(self.games, { id: _.get(self.parameters, "game.id") });
            return parameters;
        });
    };

    /**
     * Update the game parameters
     * @return {Promise}
     */
    this.saveParameters = function () {
        this.parameterBusy = true;
        return Parameters.update(null, {
            game: _.get(this.parameters, "tmpGame.id")
        }).$promise.then(function (result) {
            toaster.pop({ type: "success", title: "Game Parameters updated"});
            self.parameters.game = self.tmpGame;
            return result;
        })["finally"](function() {
            self.parameterBusy = false;
        });
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        $scope.$on("refresh-game-list", function () {
            self.loadGames();
        });
        return this.loadGames().then(function () {
            self.loadGameParameters();
        });
    };

});